package concurrency_test

import (
	"fmt"
	"github.com/benz9527/toy-box/toys/pkg/concurrency"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPollSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Poll Suite",
		types.SuiteConfig{
			LabelFilter:     "Poller",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

var _ = ginkgo.Describe("Poller Unit Tests", ginkgo.Ordered, ginkgo.Serial, func() {
	ginkgo.It("common poll",
		ginkgo.Label("Poller"),
		func() {
			invocations := 0
			taskFn := concurrency.PollTaskFunc(func() (done bool, err error) {
				invocations++
				ginkgo.GinkgoWriter.Printf("invocations: %d\n", invocations)
				return true, nil
			})
			err := concurrency.NewPoller(
				concurrency.WithPollerInterval(50*time.Millisecond),
				concurrency.WithPollerTimeout(100*time.Millisecond),
			).Poll(taskFn)
			assert.NoError(ginkgo.GinkgoT(), err)
			assert.Equal(ginkgo.GinkgoT(), 1, invocations)
		},
	)
	ginkgo.It("common poll error",
		ginkgo.Label("Poller"),
		func() {
			milliseconds := 0
			startTs := time.Now()
			taskFn := concurrency.PollTaskFunc(func() (done bool, err error) {
				milliseconds = int(time.Since(startTs).Milliseconds())
				ginkgo.GinkgoWriter.Printf("milliseconds: %d\n", milliseconds)
				return false, fmt.Errorf("test error")
			})
			err := concurrency.NewPoller(
				concurrency.WithPollerInterval(50*time.Millisecond),
				concurrency.WithPollerTimeout(100*time.Millisecond),
			).Poll(taskFn)
			assert.Error(ginkgo.GinkgoT(), err)
			assert.GreaterOrEqual(ginkgo.GinkgoT(), milliseconds, 50)
		},
	)
	ginkgo.It("poll immediately",
		ginkgo.Label("Poller"),
		func() {
			invocations := 0
			taskFn := concurrency.PollTaskFunc(func() (done bool, err error) {
				invocations++
				ginkgo.GinkgoWriter.Printf("invocations: %d\n", invocations)
				return true, nil
			})
			err := concurrency.NewPoller().PollImmediateInfinite(taskFn)
			assert.NoError(ginkgo.GinkgoT(), err)
			assert.Equal(ginkgo.GinkgoT(), 1, invocations)
		},
	)
	ginkgo.It("poll immediately error",
		ginkgo.Label("Poller"),
		func() {
			invocations := 0
			taskFn := concurrency.PollTaskFunc(func() (done bool, err error) {
				return false, fmt.Errorf("test error")
			})
			err := concurrency.NewPoller().PollImmediateInfinite(taskFn)
			assert.Error(ginkgo.GinkgoT(), err)
			assert.Equal(ginkgo.GinkgoT(), 0, invocations)
		},
	)
	ginkgo.It("poll infinite",
		ginkgo.Label("Poller"),
		func() {
			const infiniteTestTimeout = 30 * time.Second
			condC := concurrency.NewSafeChannel[struct{}]()
			doneC := concurrency.NewSafeChannel[struct{}](1)
			completedC := concurrency.NewSafeChannel[struct{}]()
			defer func() {
				_ = condC.Close()
				_ = doneC.Close()
				_ = completedC.Close()
			}()

			go func() {
				taskFn := concurrency.PollTaskFunc(func() (done bool, err error) {
					err = condC.Send(struct{}{})
					assert.NoError(ginkgo.GinkgoT(), err)
					select {
					case <-doneC.Wait():
						return true, nil
					default:

					}
					return false, nil
				})
				err := concurrency.NewPoller(
					concurrency.WithPollerInterval(time.Millisecond),
				).PollInfinite(taskFn)
				assert.NoError(ginkgo.GinkgoT(), err)
				err = condC.Close()
				assert.NoError(ginkgo.GinkgoT(), err)
				err = completedC.Send(struct{}{})
				assert.NoError(ginkgo.GinkgoT(), err)
			}()

			<-condC.Wait()

			timeoutC := time.After(infiniteTestTimeout)
			for i := 0; i < 10; i++ {
				select {
				case _, open := <-condC.Wait():
					assert.True(ginkgo.GinkgoT(), open)
					ginkgo.GinkgoWriter.Printf("received data from the channel\n")
				case <-timeoutC:
					ginkgo.GinkgoWriter.Printf("without any data caught from the channel and timeout\n")
				}
			}

			err := doneC.Send(struct{}{})
			assert.NoError(ginkgo.GinkgoT(), err)
			go func() {
				for i := 0; i < 2; i++ {
					_, open := <-condC.Wait()
					if !open {
						// Finished the whole test
						return
					}
				}
			}()
			// Blocking the whole test, unable to exit
			<-completedC.Wait()
		},
	)
	ginkgo.It("poll until",
		ginkgo.Label("Poller"),
		func() {
			// stop going to try poll channel
			stopC := concurrency.NewSafeChannel[struct{}]()
			calledC := concurrency.NewSafeChannel[bool]()
			pollDoneC := concurrency.NewSafeChannel[struct{}]()
			go func() {
				taskFn := concurrency.PollTaskFunc(func() (done bool, err error) {
					err = calledC.Send(true)
					assert.NoError(ginkgo.GinkgoT(), err)
					return false, nil
				})
				err := concurrency.NewPoller(
					concurrency.WithPollerInterval(time.Millisecond),
				).PollUntil(taskFn, stopC)
				assert.Error(ginkgo.GinkgoT(), err)
				ginkgo.GinkgoWriter.Printf("poll until error: %v\n", err)
				err = pollDoneC.Close()
				assert.NoError(ginkgo.GinkgoT(), err)
			}()

			go func() {
				for i := 0; i < 5; i++ {
					<-calledC.Wait()
					ginkgo.GinkgoWriter.Printf("received error call data from the channel\n")
				}
				err := stopC.Close()
				assert.NoError(ginkgo.GinkgoT(), err)
			}()

			<-pollDoneC.Wait()
		},
	)
})
