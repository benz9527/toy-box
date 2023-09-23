package lock

import (
	"context"
	"fmt"
	CoordinationV1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	CoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"log/slog"
	"time"
)

type podLeader struct {
	coreClient       CoreV1.CoreV1Interface
	coordinateClient CoordinationV1.CoordinationV1Interface
}

// 这个租约锁通常是用来进行 pod 的选举
// https://github.com/kubernetes/client-go/blob/master/examples/leader-election/main.go

func (p *podLeader) Election(ctx context.Context, id, namespace, name string, resourceRecords ...resourcelock.EventRecorder) error {

	_leaderID := ""
	cb := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			slog.Info(fmt.Sprintf("current pod id: %s, the leading leader id: %s", id, _leaderID))
		},
		OnStoppedLeading: func() {
			slog.Info(fmt.Sprintf("current pod id: %s, the leader id: %s no longer serve as leader", id, _leaderID))
			// os.Exit(0) // 程序整体退出
		},
		OnNewLeader: func(leaderID string) {
			_leaderID = leaderID
			slog.Info(fmt.Sprintf("try to serve as new pod leader leaderID: %s", leaderID))
		},
	}
	// 检查 context 是否被关闭
	if err := ctx.Err(); err != nil {
		slog.Warn("ctx cancelled, exit...")
		return err
	}

	cfg := resourcelock.ResourceLockConfig{
		Identity: id,
	}
	if len(resourceRecords) > 0 {
		cfg.EventRecorder = resourceRecords[0]
	}
	resLock, err := resourcelock.New(
		resourcelock.LeasesResourceLock,
		namespace, name,
		p.coreClient, p.coordinateClient,
		cfg,
	)
	if err != nil {
		return err
	}

	elector, err := leaderelection.NewLeaderElector(
		leaderelection.LeaderElectionConfig{
			Callbacks:     cb,
			Lock:          resLock,
			RenewDeadline: 2 * time.Second,
			LeaseDuration: 4 * time.Second,
			RetryPeriod:   2,
		},
	)
	if err != nil {
		return err
	}
	elector.Run(ctx)
	return nil
}
