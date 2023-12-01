//go:build !race

package concurrency_test

import (
	"math/rand"
	"runtime"
	"testing"
	"time"

	. "github.com/benz9527/toy-box/toys/pkg/concurrency"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	. "github.com/onsi/gomega"
)

func TestNewGenericChannel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NewSafeChannel Suite",
		types.SuiteConfig{
			LabelFilter:     "NewSafeChannel",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
			RandomSeed:      GinkgoRandomSeed(),
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

var _ = Describe("NewSafeChannel Unit Test",
	Label("NewSafeChannel"),
	func() {
		It("Int channel", func() {
			count := 50
			ch := NewSafeChannel[int]()
			go func(ch GenericChannel[int]) {
				for i := 0; i < count; i++ {
					time.Sleep(20 * time.Millisecond)
					if err := ch.Send(i); err != nil {
						GinkgoWriter.Printf("send data (%d) error: %v\n", i, err)
						break
					}
				}
				runtime.Goexit()
			}(ch)
			go func(ch GenericChannel[int]) {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				n := r.Intn(count * 20)
				time.Sleep(time.Duration(n) * time.Millisecond)
				_ = ch.Close()
				GinkgoWriter.Printf("close channel after %d ms\n", n)
				runtime.Goexit()
			}(ch)
			for {
				select {
				case v, ok := <-ch.Wait():
					if !ok {
						GinkgoWriter.Printf("channel closed\n")
						time.Sleep(50 * time.Millisecond)
						return
					}
					GinkgoWriter.Printf("receive data (%d)\n", v)
				}
			}
		})
		It("Int buffer channel", func() {
			count := 50
			ch := NewSafeChannel[int](10)
			go func(ch GenericChannel[int]) {
				for i := 0; i < count; i++ {
					time.Sleep(20 * time.Millisecond)
					if err := ch.Send(i); err != nil {
						GinkgoWriter.Printf("send data (%d) to buffer channel error: %v\n", i, err)
						break
					}
				}
				runtime.Goexit()
			}(ch)
			go func(ch GenericChannel[int]) {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				n := r.Intn(count * 20)
				time.Sleep(time.Duration(n) * time.Millisecond)
				_ = ch.Close()
				GinkgoWriter.Printf("close buffer channel after %d ms\n", n)
				runtime.Goexit()
			}(ch)
			for {
				select {
				case v, ok := <-ch.Wait():
					if !ok {
						GinkgoWriter.Printf("buffer channel closed\n")
						time.Sleep(50 * time.Millisecond)
						return
					}
					GinkgoWriter.Printf("receive data (%d)\n", v)
				}
			}
		})
	},
)
