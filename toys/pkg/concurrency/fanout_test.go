package concurrency_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func fanOutWorker(wg *sync.WaitGroup, id int, jobs <-chan int, results chan<- int) {
	defer wg.Done()
	for job := range jobs {
		GinkgoWriter.Printf("worker %d started job %d\n", id, job)
		time.Sleep(time.Second)
		results <- job * 2
	}
}

var _ = Describe("Easy fanOut concurrency test", func() {
	wg := &sync.WaitGroup{}
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	ans := 0
	It("should fan out jobs to workers", func(ctx SpecContext) {
		wg.Add(3)
		for i := 1; i <= 3; i++ {
			go fanOutWorker(wg, i, jobs, results)
		}
		go func() {
			<-ctx.Done()
			By("context timeout, unable to close jobs and results")
		}()
		go func() {
			wg.Wait()
			close(results)
		}()
		for i := 1; i <= 5; i++ {
			jobs <- i
		}
		close(jobs)
		for result := range results {
			ans += result
		}
		Expect(ans).To(Equal(30))
	}, SpecTimeout(3*time.Second))
})
