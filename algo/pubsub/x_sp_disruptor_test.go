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

func testXSinglePipelineDisruptor(t *testing.T, gTotal, tasks int, bs BlockStrategy) {
	counter := &atomic.Int64{}
	wg := &sync.WaitGroup{}
	wg.Add(gTotal)
	disruptor := NewXSinglePipelineDisruptor[int](1024*1024,
		bs,
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
		bs     BlockStrategy
	}{
		{10, 100, NewXGoSchedBlockStrategy()},
		{100, 10000, NewXGoSchedBlockStrategy()},
		{500, 10000, NewXGoSchedBlockStrategy()},
		{1000, 10000, NewXGoSchedBlockStrategy()},
		{5000, 10000, NewXGoSchedBlockStrategy()},
		{10000, 10000, NewXGoSchedBlockStrategy()},
		{5000, 10000, NewXNoCacheChannelBlockStrategy()},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("gTotal: %d, tasks: %d", tc.gTotal, tc.tasks), func(t *testing.T) {
			testXSinglePipelineDisruptor(t, tc.gTotal, tc.tasks, tc.bs)
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

func testXSinglePipelineDisruptorWithRandomSleep(t *testing.T, num, capacity int) {
	wg := &sync.WaitGroup{}
	wg.Add(num)
	results := map[string]struct{}{}
	disruptor := NewXSinglePipelineDisruptor[string](uint64(capacity),
		NewXNoCacheChannelBlockStrategy(),
		func(event string) error {
			nextInt := rand.Intn(100)
			time.Sleep(time.Duration(nextInt) * time.Millisecond)
			results[event] = struct{}{}
			wg.Done()
			return nil
		},
	)
	if err := disruptor.Start(); err != nil {
		t.Fatalf("disruptor start failed, err: %v", err)
	}
	for i := 0; i < num; i++ {
		if _, _, err := disruptor.Publish(fmt.Sprintf("event-%d", i)); err != nil {
			t.Logf("publish failed, err: %v", err)
		}
	}
	wg.Wait()
	err := disruptor.Stop()
	assert.NoError(t, err)
	assert.Equal(t, num, len(results))
	for i := 0; i < num; i++ {
		assert.Contains(t, results, fmt.Sprintf("event-%d", i))
	}
}

func TestXSinglePipelineDisruptorWithRandomSleepEvent(t *testing.T) {
	testcases := []struct {
		num      int
		capacity int
	}{
		{10, 2},
		{100, 4},
		{200, 10},
		{500, 20},
	}
	loops := 2
	for i := 0; i < loops; i++ {
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("num: %d, capacity: %d", tc.num, tc.capacity), func(t *testing.T) {
				testXSinglePipelineDisruptorWithRandomSleep(t, tc.num, tc.capacity)
			})
		}
	}
}

func TestXSinglePipelineDisruptor_PublishTimeout(t *testing.T) {
	num := 10
	disruptor := NewXSinglePipelineDisruptor[string](2,
		NewXGoSchedBlockStrategy(),
		func(event string) error {
			nextInt := rand.Intn(10)
			if nextInt == 0 {
				nextInt = 2
			}
			time.Sleep(time.Duration(nextInt) * time.Millisecond)
			slog.Info("handle event details", "name", event)
			return nil
		},
	)
	if err := disruptor.Start(); err != nil {
		t.Fatalf("disruptor start failed, err: %v", err)
	}
	for i := 0; i < num; i++ {
		event := fmt.Sprintf("event-%d", i)
		disruptor.PublishTimeout(event, 5*time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)
	err := disruptor.Stop()
	assert.NoError(t, err)
}
