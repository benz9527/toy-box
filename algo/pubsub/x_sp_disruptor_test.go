package pubsub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCeilCapacity(t *testing.T) {
	testcases := []struct {
		capacity uint64
		ceil     uint64
	}{
		{0, 2},
		{1, 2},
		{2, 2},
		{3, 4},
		{4, 4},
		{7, 8},
		{8, 8},
		{9, 16},
		{16, 16},
		{31, 32},
		{32, 32},
		{58, 64},
		{64, 64},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("capacity: %d, ceil: %d", tc.capacity, tc.ceil), func(t *testing.T) {
			assert.Equal(t, tc.ceil, ceilCapacity(tc.capacity))
		})
	}
}

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

func TestXSinglePipelineDisruptorWithRandomSleepEventHandler(t *testing.T) {
	num := 10
	wg := &sync.WaitGroup{}
	wg.Add(num)
	disruptor := NewXSinglePipelineDisruptor[string](2,
		NewXGoSchedBlockStrategy(),
		func(event string) error {
			nextInt := rand.Intn(100)
			time.Sleep(time.Duration(nextInt) * time.Millisecond)
			slog.Info("event details", "name", event)
			wg.Done()
			return nil
		},
	)
	if err := disruptor.Start(); err != nil {
		t.Fatalf("disruptor start failed, err: %v", err)
	}
	for i := 0; i < 10; i++ {
		if _, _, err := disruptor.Publish(fmt.Sprintf("event-%d", i)); err != nil {
			t.Logf("publish failed, err: %v", err)
			break
		}
	}
	wg.Wait()
	err := disruptor.Stop()
	assert.NoError(t, err)
}
