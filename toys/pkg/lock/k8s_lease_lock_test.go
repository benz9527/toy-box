package lock

import (
	"context"
	"errors"
	"fmt"
	Ants "github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	CoordinationV1Meta "k8s.io/api/coordination/v1"
	KErr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	KTesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

// https://github.com/kubernetes/client-go/blob/master/examples/fake-client/main_test.go
// https://github.com/kubernetes/client-go/blob/master/tools/leaderelection/leaderelection_test.go

type reactor struct {
	verb     string
	objType  string
	reaction KTesting.ReactionFunc
}

func TestK8SCoordinateLeaseLockByFakeClient(t *testing.T) {

	// acquire from no object
	var globalLeaseLock runtime.Object
	reactors := []reactor{
		{
			verb:    "get",
			objType: "leases",
			reaction: func(action KTesting.Action) (handled bool, ret runtime.Object, err error) {
				if globalLeaseLock == nil {
					return true, nil, KErr.NewNotFound(
						action.(KTesting.GetAction).GetResource().GroupResource(),
						action.(KTesting.GetAction).GetName(),
					)
				}
				return true, globalLeaseLock, nil
			},
		}, {
			verb:    "create",
			objType: "leases",
			reaction: func(action KTesting.Action) (handled bool, ret runtime.Object, err error) {
				if globalLeaseLock == nil {
					globalLeaseLock = action.(KTesting.CreateAction).GetObject()
					return true, globalLeaseLock, nil
				}
				return true, nil, KErr.NewAlreadyExists(
					action.(KTesting.CreateAction).GetResource().GroupResource(),
					"mylease",
				)
			},
		}, {
			verb:    "update",
			objType: "leases",
			reaction: func(action KTesting.Action) (handled bool, ret runtime.Object, err error) {
				if globalLeaseLock != nil {
					globalLeaseLock = action.(KTesting.UpdateAction).GetObject()
				}
				return true, globalLeaseLock, nil
			},
		},
	}

	cs := fake.NewSimpleClientset() // 需要去掉默认的 reactor 行为才能调用上面定义的 reactor 行为
	cs.ReactionChain = make([]KTesting.Reactor, 0, 8)
	// 记录配对的操作
	watcherStarted := make(chan struct{})
	cs.PrependWatchReactor("*", func(action KTesting.Action) (handled bool, ret watch.Interface, err error) {
		defer close(watcherStarted)
		gvr := action.GetResource()
		namespace := action.GetNamespace()
		_watch, err := cs.Tracker().Watch(gvr, namespace)
		if err != nil {
			return false, nil, err
		}
		return true, _watch, nil
	})
	for _, reactor := range reactors {
		cs.AddReactor(reactor.verb, reactor.objType, reactor.reaction)
	}

	ctx, cancel := context.WithTimeoutCause(context.Background(), 10*time.Second, errors.New("lease test end of life"))
	defer func() {
		cancel()
	}()

	// 这个 cache informer 有点鸡肋
	_informers := informers.NewSharedInformerFactory(cs, 0)
	leaseInformer := _informers.Coordination().V1().Leases().Informer()
	_, _ = leaseInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			lease, ok := obj.(*CoordinationV1Meta.Lease)
			assert.True(t, ok)
			slog.Info(fmt.Sprintf("coordinate lease id: %d", lease.Spec.HolderIdentity))
		},
		DeleteFunc: func(obj interface{}) {
			// 无效过程，暂时无法模拟删除租约锁，再被抢占和重新创建的过程
			lease, ok := obj.(*CoordinationV1Meta.Lease)
			assert.True(t, ok)
			gLease, ok := globalLeaseLock.(*CoordinationV1Meta.Lease)
			assert.True(t, ok)
			slog.Info("deleting...")
			if lease.Spec.HolderIdentity == gLease.Spec.HolderIdentity &&
				lease.GetName() == gLease.GetName() &&
				lease.GetNamespace() == gLease.GetNamespace() {
				globalLeaseLock = nil
			}
		},
	})
	_informers.Start(ctx.Done())
	cache.WaitForCacheSync(ctx.Done(), leaseInformer.HasSynced)
	<-watcherStarted

	leader := podLeader{
		coreClient:       cs.CoreV1(),
		coordinateClient: cs.CoordinationV1(),
	}
	hn, err := os.Hostname()
	assert.NoError(t, err)
	assert.NotEqual(t, "", hn)
	maxCount := 10
	wg := sync.WaitGroup{}
	wg.Add(maxCount)
	rnd := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < maxCount; i++ {
		_ = Ants.Submit(func(idx int) func() {
			return func() {
				time.Sleep(time.Duration(rnd.Intn(maxCount)-1) * time.Second)
				defer wg.Done()
				slog.Info(fmt.Sprintf("my goroutine id: %d", idx))
				_ = leader.Election(ctx, fmt.Sprintf("%s-%d", hn, idx), "ben-lease", "mylease", record.NewFakeRecorder(100))
			}
		}(i))
	}
	wg.Wait()
}
