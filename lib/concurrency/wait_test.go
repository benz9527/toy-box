//go:build !race

package concurrency_test

import (
	"context"
	"errors"
	"github.com/benz9527/toy-box/lib/concurrency"
	"github.com/benz9527/toy-box/lib/runtime"
	"github.com/onsi/ginkgo/v2"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestWaitJitterSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Wait Jitter Suite",
		types.SuiteConfig{
			LabelFilter:     "Jitter",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

var _ = ginkgo.Describe("Jitter Unit Tests", ginkgo.Ordered, ginkgo.Serial, func() {
	ginkgo.It("jitter recover from panic",
		ginkgo.Label("Jitter"), func() {
			called := 0
			handled := 0
			stopC := concurrency.NewSafeChannel[struct{}]()
			resetFn := runtime.SetDefaultCrashHandlers(runtime.NoOpCrashHandler)
			defer func() {
				resetFn()
			}()

			concurrency.NewJitter(
				time.Millisecond, 1.0,
				concurrency.WithJitterSliding(true),
				concurrency.WithJitterStopChannel(stopC),
				concurrency.WithJitterCrashHandlers(
					func(any) {
						handled++
					},
					func(any) {
						ginkgo.GinkgoWriter.Printf("recover from panic\n")
					},
				),
			).Until(func() {
				if called > 2 {
					_ = stopC.Close()
					return
				}
				called++
				panic(errors.New("test jitter until recover from panic"))
			})
			assert.Equal(ginkgo.GinkgoT(), 3, handled)
		},
	)
	ginkgo.It("jitter until with context",
		ginkgo.Label("Jitter"), func() {
			ctx, cancel := context.WithCancel(context.TODO())
			calledC := concurrency.NewSafeChannel[struct{}]()
			go func() {
				concurrency.NewJitter(0, concurrency.Factor1x,
					concurrency.WithJitterSliding(true),
					concurrency.WithJitterTraceID("ctx-util-1")).
					UntilWithContext(ctx, func(ctx context.Context) {
						_ = calledC.Send(struct{}{})
						ginkgo.GinkgoWriter.Printf("context cancel called in util\n")
					})
			}()
			<-calledC.Wait()
			cancel()
			_ = calledC.Close()
			<-calledC.Wait() // Parallel may be blocked here.
		},
	)
	ginkgo.It("jitter until non sliding",
		ginkgo.Label("Jitter"), func() {
			stopC := concurrency.NewSafeChannel[struct{}]()
			_ = stopC.Close()
			concurrency.NewJitter(0, concurrency.Factor1x,
				concurrency.WithJitterSliding(true),
				concurrency.WithJitterStopChannel(stopC)).
				NonSlidingUntil(func() {
					ginkgo.GinkgoWriter.Printf("jitter until non sliding should not have been invoked\n")
				})

			stopC = concurrency.NewSafeChannel[struct{}]()
			calledC := concurrency.NewSafeChannel[struct{}]()
			go func() {
				concurrency.NewJitter(0, concurrency.Factor1x,
					concurrency.WithJitterSliding(true),
					concurrency.WithJitterStopChannel(stopC)).
					NonSlidingUntil(func() {
						_ = calledC.Send(struct{}{})
						ginkgo.GinkgoWriter.Printf("called non sliding\n")
					})
			}()
			<-calledC.Wait()
			_ = stopC.Close()
			_ = calledC.Close()
			<-calledC.Wait() // Parallel may be blocked here.
		},
	)
	ginkgo.It("jitter until with context non sliding",
		ginkgo.Label("Jitter"), func() {
			ctx, cancel := context.WithCancel(context.TODO())
			calledC := concurrency.NewSafeChannel[struct{}]()
			go func() {
				concurrency.NewJitter(0, concurrency.Factor1x,
					concurrency.WithJitterTraceID("ctx-util-no-sliding-1")).
					NonSlidingUntilWithContext(ctx, func(ctx context.Context) {
						_ = calledC.Send(struct{}{})
						ginkgo.GinkgoWriter.Printf("context cancel called in until non sliding\n")
					})
			}()
			<-calledC.Wait()
			cancel()
			_ = calledC.Close()
			<-calledC.Wait() // Parallel may be blocked here.
		},
	)
	ginkgo.It("jitter until returns immediately",
		ginkgo.Label("Jitter"), func() {
			startTs := time.Now()
			stopC := concurrency.NewSafeChannel[struct{}]()
			concurrency.NewJitter(30*time.Second, concurrency.Factor1x,
				concurrency.WithJitterSliding(true),
				concurrency.WithJitterStopChannel(stopC)).
				Until(func() {
					ginkgo.GinkgoWriter.Printf("jitter until returns immediately close channel\n")
					_ = stopC.Close()
				})
			assert.Falsef(ginkgo.GinkgoT(), startTs.Add(25*time.Second).Before(time.Now()), "jitter until returns immediately")
		},
	)
	ginkgo.It("jitter until returns immediately factor 0",
		ginkgo.Label("Jitter"), func() {
			startTs := time.Now()
			stopC := concurrency.NewSafeChannel[struct{}]()
			concurrency.NewJitter(30*time.Second, concurrency.Factor0x,
				concurrency.WithJitterSliding(true),
				concurrency.WithJitterStopChannel(stopC)).
				Until(func() {
					ginkgo.GinkgoWriter.Printf("jitter until returns immediately close channel\n")
					_ = stopC.Close()
				})
			assert.Falsef(ginkgo.GinkgoT(), startTs.Add(25*time.Second).Before(time.Now()), "jitter until returns immediately")
		},
	)
	ginkgo.It("jitter until with negative factor",
		ginkgo.Label("Jitter"), func() {
			startTs := time.Now()
			stopC := concurrency.NewSafeChannel[struct{}]()
			receivedC := make(chan struct{})
			calledC := concurrency.NewSafeChannel[struct{}]()

			go func() {
				concurrency.NewJitter(time.Second, -30.0,
					concurrency.WithJitterSliding(true),
					concurrency.WithJitterStopChannel(stopC)).
					Until(func() {
						_ = calledC.Send(struct{}{})
						<-receivedC
						ginkgo.GinkgoWriter.Printf("jitter until with negative factor received\n")
					})
			}()

			// 1st loop
			<-calledC.Wait()
			receivedC <- struct{}{}

			// 2nd loop
			<-calledC.Wait()
			_ = stopC.Close()
			receivedC <- struct{}{}

			assert.Falsef(ginkgo.GinkgoT(), startTs.Add(3*time.Second).Before(time.Now()),
				"jitter until with negative factor did not returned after predefined period when then stop chan was closed inside the func")
		},
	)
})
