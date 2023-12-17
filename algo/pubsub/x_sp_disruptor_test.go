package pubsub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testXSinglePipelineDisruptor(t *testing.T, gTotal, tasks int) {
	counter := &atomic.Int64{}
	wg := &sync.WaitGroup{}
	wg.Add(gTotal)
	disruptor := NewXSinglePipelineDisruptor[int](1024*1024,
		NewXGoSchedBlockStrategy(),
		func(event int) error {
			counter.Add(1)
			return nil
		},
	)
	if err := disruptor.Start(); err != nil {
		t.Fatalf("disruptor start failed, err: %v", err)
	}
	beginTs := time.Now()
	for i := 0; i < gTotal; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < tasks; j++ {
				if _, _, err := disruptor.Publish(j); err != nil {
					t.Logf("publish failed, err: %v", err)
					break
				}
			}
		}()
	}
	wg.Wait()
	diff := time.Now().Sub(beginTs)
	t.Logf("total: %d, tasks: %d, cost: %v, tps: %v/s", gTotal, tasks, diff, float64(gTotal*tasks)/diff.Seconds())
	time.Sleep(time.Second)
	assert.Equal(t, int64(gTotal*tasks), counter.Load())
	err := disruptor.Stop()
	assert.NoError(t, err)
}

func TestXSinglePipelineDisruptor(t *testing.T) {
	testcases := []struct {
		gTotal int
		tasks  int
	}{
		{100, 10000},
		{500, 10000},
		{1000, 10000},
		{5000, 10000},
		{10000, 10000},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("gTotal: %d, tasks: %d", tc.gTotal, tc.tasks), func(t *testing.T) {
			testXSinglePipelineDisruptor(t, tc.gTotal, tc.tasks)
		})
	}
}

func testNoCacheChannel(t *testing.T, chSize, gTotal, tasks int) {
	counter := &atomic.Int64{}
	wg := &sync.WaitGroup{}
	wg.Add(gTotal)
	var ch chan int
	if chSize > 0 {
		ch = make(chan int, chSize)
	} else {
		ch = make(chan int)
	}
	go func() {
		for range ch {
			counter.Add(1)
		}
	}()
	beginTs := time.Now()
	for i := 0; i < gTotal; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < tasks; j++ {
				ch <- j
			}
		}()
	}
	wg.Wait()
	diff := time.Now().Sub(beginTs)
	t.Logf("total: %d, tasks: %d, cost: %v, tps: %v/s", gTotal, tasks, diff, float64(gTotal*tasks)/diff.Seconds())
	time.Sleep(time.Second)
	assert.Equal(t, int64(gTotal*tasks), counter.Load())
}

func TestNoCacheChannel(t *testing.T) {
	testcases := []struct {
		gTotal int
		tasks  int
	}{
		{100, 10000},
		{500, 10000},
		{1000, 10000},
		{5000, 10000},
		{10000, 10000},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("gTotal: %d, tasks: %d", tc.gTotal, tc.tasks), func(t *testing.T) {
			testNoCacheChannel(t, 0, tc.gTotal, tc.tasks)
		})
	}
}

func TestCacheChannel(t *testing.T) {
	testcases := []struct {
		gTotal int
		tasks  int
	}{
		{100, 10000},
		{500, 10000},
		{1000, 10000},
		{5000, 10000},
		{10000, 10000},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("gTotal: %d, tasks: %d", tc.gTotal, tc.tasks), func(t *testing.T) {
			testNoCacheChannel(t, 1024*1024, tc.gTotal, tc.tasks)
		})
	}
}
