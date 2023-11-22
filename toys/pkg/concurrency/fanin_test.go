package concurrency_test

import (
	"github.com/onsi/ginkgo/v2/types"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFanInSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FanIn Suite",
		types.SuiteConfig{
			LabelFilter:     "FanIn",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

func fanInProducer(C chan<- int, num int) {
	for i := 0; i < 5; i++ {
		C <- num + i
	}
	close(C)
}

func fanIn(doneC <-chan struct{}, C ...<-chan int) <-chan int {
	outC := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(len(C))
	// 这种 fanIn 不能动态聚合通道，只能在初始化时确定聚合的通道
	// 不过通道有缓冲区间，完全可以不需要多通道聚合
	for _, c := range C {
		go func(c <-chan int) {
			defer wg.Done()
			for dataC := range c {
				// 必须卡在这个等待数据，使用 default 会导致数据丢失
				select {
				case outC <- dataC:
				case <-doneC:
					return
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(outC)
	}()

	return outC
}

var _ = Describe("Easy fan in concurrency test", Label("FanIn"),
	func() {
		C1, C2, doneC := make(chan int), make(chan int), make(chan struct{})
		ans := 0
		It("should fan in jobs from producers", func(ctx SpecContext) {
			go fanInProducer(C1, 0)
			go fanInProducer(C2, 10)
			for v := range fanIn(doneC, C1, C2) {
				GinkgoWriter.Printf("v: %d\n", v)
				ans += v
			}
			Expect(ans).To(Equal(70))
		},
			SpecTimeout(3*time.Second))
	},
)
