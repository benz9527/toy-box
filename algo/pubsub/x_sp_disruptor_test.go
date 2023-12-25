package pubsub

import (
	"bytes"
	"fmt"
	"github.com/benz9527/toy-box/algo/bitmap"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
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

func testXSinglePipelineDisruptorUint64(
	t *testing.T, gTotal, tasks int, capacity uint64,
	bs BlockStrategy, bitmapCheck bool,
	reportFile *os.File, errorCounter *atomic.Uint64,
) {
	var (
		counter = &atomic.Int64{}
		bm      bitmap.Bitmap
		checkBM bitmap.Bitmap
	)
	checkBM = bitmap.NewX32Bitmap(uint64(gTotal * tasks))
	if bitmapCheck {
		bm = bitmap.NewX32Bitmap(uint64(gTotal * tasks))
	}
	wg := &sync.WaitGroup{}
	wg.Add(gTotal)
	rwg := &sync.WaitGroup{}
	rwg.Add(gTotal * tasks)
	disruptor := NewXSinglePipelineDisruptor[uint64](capacity,
		bs,
		func(event uint64) error {
			defer rwg.Done()
			counter.Add(1)
			if bitmapCheck {
				bm.SetBit(event, true)
			}
			return nil
		},
	)
	if err := disruptor.Start(); err != nil {
		t.Fatalf("disruptor start failed, err: %v", err)
	}
	for i := 0; i < gTotal; i++ {
		for j := 0; j < tasks; j++ {
			checkBM.SetBit(uint64(i*tasks+j), true)
		}
	}
	beginTs := time.Now()
	for i := 0; i < gTotal; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < tasks; j++ {
				if _, _, err := disruptor.Publish(uint64(idx*tasks + j)); err != nil {
					t.Logf("publish failed, err: %v", err)
					if errorCounter != nil {
						errorCounter.Add(1)
					}
					break
				}
			}
		}(i)
	}
	wg.Wait()
	diff := time.Now().Sub(beginTs)
	summary := fmt.Sprintf("published total: %d, tasks: %d, cost: %v, tps: %v/s", gTotal, tasks, diff, float64(gTotal*tasks)/diff.Seconds())
	t.Log(summary)
	if reportFile != nil {
		_, _ = reportFile.WriteString(summary + "\n")
	}
	rwg.Wait()
	if reportFile == nil {
		time.Sleep(time.Second)
		assert.Equal(t, int64(gTotal*tasks), counter.Load())
	}
	err := disruptor.Stop()
	assert.NoError(t, err)
	if bitmapCheck {
		if reportFile != nil {
			_, _ = reportFile.WriteString(fmt.Sprintf("gTotal(%d), tasks(%d):\n", gTotal, tasks))
		}
		bm1bits := bm.GetBits()
		bm2bits := checkBM.GetBits()
		if !bm.EqualTo(checkBM) {
			if reportFile != nil {
				_, _ = reportFile.WriteString("bitmap check failed by not equal!\n")
			}
			if errorCounter != nil {
				errorCounter.Add(1)
			}
			for i := 0; i < len(bm1bits); i++ {
				if bytes.Compare(bm1bits[i:i+1], bm2bits[i:i+1]) != 0 {
					if reportFile != nil {
						_, _ = reportFile.WriteString(fmt.Sprintf("idx: %d, bm1: %08b, bm2: %08b\n", i, bm1bits[i:i+1], bm2bits[i:i+1]))
					}
					t.Logf("idx: %d, bm1: %08b, bm2: %08b\n", i, bm1bits[i:i+1], bm2bits[i:i+1])
				}
			}
		}
		// check store whether contains zero bits
		if reportFile != nil {
			_, _ = reportFile.WriteString("check store whether contains zero bits(exclude the last one):\n")
			for i := 0; i < len(bm2bits)-1; i++ {
				if bm2bits[i]&0xf != 0xf {
					_, _ = reportFile.WriteString(fmt.Sprintf("store idx: %d, bm2: %08b\n", i, bm2bits[i:i+1]))
				}
			}
			_, _ = reportFile.WriteString("====== end report ======\n")
		}
	}
	if bm != nil {
		bm.Free()
	}
	if checkBM != nil {
		checkBM.Free()
	}
}

func testXSinglePipelineDisruptorString(t *testing.T, gTotal, tasks int, capacity uint64, bs BlockStrategy, bitmapCheck bool, reportFile *os.File, errorCounter *atomic.Uint64) {
	var (
		counter = &atomic.Int64{}
		bm      bitmap.Bitmap
		checkBM bitmap.Bitmap
	)
	if bitmapCheck {
		bm = bitmap.NewX32Bitmap(uint64(gTotal * tasks))
		checkBM = bitmap.NewX32Bitmap(uint64(gTotal * tasks))
	}
	wg := &sync.WaitGroup{}
	wg.Add(gTotal)
	rwg := &sync.WaitGroup{}
	rwg.Add(gTotal * tasks)
	disruptor := NewXSinglePipelineDisruptor[string](capacity,
		bs,
		func(event string) error {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("error panic: %v", r)
					if reportFile != nil {
						_, _ = reportFile.WriteString(fmt.Sprintf("error panic: %v\n", r))
					}
					if errorCounter != nil {
						errorCounter.Add(1)
					}
				}
				rwg.Done()
			}()
			counter.Add(1)
			if bitmapCheck {
				e, err := strconv.ParseUint(event, 10, 64)
				if err != nil {
					t.Logf("error parse uint64 failed, err: %v", err)
					if reportFile != nil {
						_, _ = reportFile.WriteString(fmt.Sprintf("error parse uint64 failed, err: %v\n", err))
					}
				}
				bm.SetBit(e, true)
			}
			if event == "" {
				t.Logf("error event is empty, counter: %d", counter.Load())
			}
			return nil
		},
	)
	if err := disruptor.Start(); err != nil {
		t.Fatalf("disruptor start failed, err: %v", err)
	}
	for i := 0; i < gTotal; i++ {
		for j := 0; j < tasks; j++ {
			checkBM.SetBit(uint64(i*tasks+j), true)
		}
	}
	beginTs := time.Now()
	for i := 0; i < gTotal; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < tasks; j++ {
				if _, _, err := disruptor.Publish(fmt.Sprintf("%d", idx*tasks+j)); err != nil {
					t.Logf("publish failed, err: %v", err)
					if errorCounter != nil {
						errorCounter.Add(1)
					}
					break
				}
			}
		}(i)
	}
	wg.Wait()
	diff := time.Now().Sub(beginTs)
	summary := fmt.Sprintf("published total: %d, tasks: %d, cost: %v, tps: %v/s", gTotal, tasks, diff, float64(gTotal*tasks)/diff.Seconds())
	t.Log(summary)
	if reportFile != nil {
		_, _ = reportFile.WriteString(summary + "\n")
	}
	rwg.Wait()
	if reportFile == nil {
		time.Sleep(time.Second)
		assert.Equal(t, int64(gTotal*tasks), counter.Load())
	}
	err := disruptor.Stop()
	assert.NoError(t, err)
	if bitmapCheck {
		if reportFile != nil {
			_, _ = reportFile.WriteString(fmt.Sprintf("gTotal(%d), tasks(%d):\n", gTotal, tasks))
		}
		bm1bits := bm.GetBits()
		bm2bits := checkBM.GetBits()
		if !bm.EqualTo(checkBM) {
			if reportFile != nil {
				_, _ = reportFile.WriteString("bitmap check failed by not equal!\n")
			}
			if errorCounter != nil {
				errorCounter.Add(1)
			}
			for i := 0; i < len(bm1bits); i++ {
				if bytes.Compare(bm1bits[i:i+1], bm2bits[i:i+1]) != 0 {
					if reportFile != nil {
						_, _ = reportFile.WriteString(fmt.Sprintf("error store idx: %d, bm1: %08b, bm2: %08b\n", i, bm1bits[i:i+1], bm2bits[i:i+1]))
					}
					t.Logf("idx: %d, bm1: %08b, bm2: %08b\n", i, bm1bits[i:i+1], bm2bits[i:i+1])
				}
			}
		}
		// check store whether contains zero bits
		if reportFile != nil {
			_, _ = reportFile.WriteString("check store whether contains zero bits(exclude the last one):\n")
			for i := 0; i < len(bm2bits)-1; i++ {
				if bm2bits[i]&0xf != 0xf {
					_, _ = reportFile.WriteString(fmt.Sprintf("error store idx: %d, bm2: %08b\n", i, bm2bits[i:i+1]))
				}
			}
			_, _ = reportFile.WriteString("====== end report ======\n")
		}
	}
	if bm != nil {
		bm.Free()
	}
	if checkBM != nil {
		checkBM.Free()
	}
}

func TestXSinglePipelineDisruptor(t *testing.T) {
	testcases := []struct {
		name   string
		gTotal int
		tasks  int
		bs     BlockStrategy
	}{
		{"gosched 10*100", 10, 100, NewXGoSchedBlockStrategy()},
		{"gosched 100*10000", 100, 10000, NewXGoSchedBlockStrategy()},
		{"gosched 500*10000", 500, 10000, NewXGoSchedBlockStrategy()},
		{"gosched 1000*10000", 1000, 10000, NewXGoSchedBlockStrategy()},
		{"gosched 5000*10000", 5000, 10000, NewXGoSchedBlockStrategy()},
		{"gosched 10000*10000", 10000, 10000, NewXGoSchedBlockStrategy()},
		{"nochan 5000*10000", 5000, 10000, NewXCacheChannelBlockStrategy()},
		{"cond 5000*10000", 5000, 10000, NewXCondBlockStrategy()},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			testXSinglePipelineDisruptorUint64(t, tc.gTotal, tc.tasks, 1024*1024, tc.bs, false, nil, nil)
		})
	}
}

func TestXSinglePipelineDisruptorWithBitmapCheck(t *testing.T) {
	testcases := []struct {
		name   string
		gTotal int
		tasks  int
		bs     BlockStrategy
	}{
		{"gosched 1*10000", 1, 10000, NewXGoSchedBlockStrategy()},
		{"nocachech 1*10000", 1, 10000, NewXCacheChannelBlockStrategy()},
		{"cond 1*10000", 1, 10000, NewXCondBlockStrategy()},
		//{"gosched 10*100", 10, 100, NewXGoSchedBlockStrategy()},
		//{"gosched 100*10000", 100, 10000, NewXGoSchedBlockStrategy()},
		//{"gosched 500*10000", 500, 10000, NewXGoSchedBlockStrategy()},
		//{"gosched 1000*10000", 1000, 10000, NewXGoSchedBlockStrategy()},
		//{"gosched 5000*10000", 5000, 10000, NewXGoSchedBlockStrategy()},
		//{"gosched 10000*10000", 10000, 10000, NewXGoSchedBlockStrategy()},
		//{"nochan 5000*10000", 5000, 10000, NewXCacheChannelBlockStrategy()},
		//{"cond 5000*10000", 5000, 10000, NewXCondBlockStrategy()},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			testXSinglePipelineDisruptorUint64(t, tc.gTotal, tc.tasks, 1024*1024, tc.bs, true, nil, nil)
		})
	}
}

func TestXSinglePipelineDisruptorWithBitmapCheckAndReport(t *testing.T) {
	errorCounter := &atomic.Uint64{}
	reportFile, err := os.OpenFile(filepath.Join(os.TempDir(), "pubsub-report-"+time.Now().Format(time.RFC3339)+".txt"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer func() {
		if reportFile != nil {
			_ = reportFile.Close()
		}
	}()
	assert.NoError(t, err)
	testcases := []struct {
		name     string
		loop     int
		gTotal   int
		tasks    int
		capacity uint64
		bs       BlockStrategy
	}{
		{"gosched 1*10000", 10000, 1, 10000, 1024, NewXGoSchedBlockStrategy()},
		{"gosched 10*100", 1000, 10, 100, 512, NewXGoSchedBlockStrategy()},
		{"gosched 10*100", 1000, 10, 100, 1024, NewXGoSchedBlockStrategy()},
		{"gosched 100*10000", 200, 100, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 500*10000", 10, 500, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 1000*10000", 10, 1000, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 5000*10000", 10, 5000, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 10000*10000", 5, 10000, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		//{"chan 1*10000", 10000, 1, 10000, 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 10*100", 1000, 10, 100, 512, NewXCacheChannelBlockStrategy()},
		//{"chan 10*100", 1000, 10, 100, 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 100*10000", 200, 100, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 500*10000", 10, 500, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 1000*10000", 10, 1000, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 5000*10000", 10, 5000, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 10000*10000", 5, 10000, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"cond 1*10000", 10000, 1, 10000, 1024, NewXCondBlockStrategy()},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.loop; i++ {
				_, _ = reportFile.WriteString(fmt.Sprintf("\n====== begin uint64 report(%s, %d) ======\n", tc.name, i))
				testXSinglePipelineDisruptorUint64(t, tc.gTotal, tc.tasks, tc.capacity, tc.bs, true, reportFile, errorCounter)
			}
		})
	}
}

func TestXSinglePipelineDisruptorWithBitmapCheckAndReport_str(t *testing.T) {
	errorCounter := &atomic.Uint64{}
	reportFile, err := os.OpenFile(filepath.Join(os.TempDir(), "pubsub-report-str-"+time.Now().Format(time.RFC3339)+".txt"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer func() {
		if reportFile != nil {
			_ = reportFile.Close()
		}
	}()
	assert.NoError(t, err)
	testcases := []struct {
		name     string
		loop     int
		gTotal   int
		tasks    int
		capacity uint64
		bs       BlockStrategy
	}{
		{"gosched 1*10000 str", 1000, 1, 10000, 1024, NewXGoSchedBlockStrategy()},
		{"gosched 10*100 str", 1000, 10, 100, 512, NewXGoSchedBlockStrategy()},
		{"gosched 10*100 str", 1000, 10, 100, 1024, NewXGoSchedBlockStrategy()},
		{"gosched 100*10000 str", 1000, 100, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 500*10000 str", 100, 500, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 1000*10000 str", 10, 1000, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 5000*10000 str", 10, 5000, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		{"gosched 10000*10000 str", 5, 10000, 10000, 1024 * 1024, NewXGoSchedBlockStrategy()},
		//{"chan 1*10000 str", 1000, 1, 10000, 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 10*100 str", 1000, 10, 100, 512, NewXCacheChannelBlockStrategy()},
		//{"chan 10*100 str", 1000, 10, 100, 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 100*10000 str", 1000, 100, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 500*10000 str", 10, 500, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 1000*10000 str", 10, 1000, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 5000*10000 str", 10, 5000, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"chan 10000*10000 str", 5, 10000, 10000, 1024 * 1024, NewXCacheChannelBlockStrategy()},
		//{"cond 1*10000 str", 1000, 1, 10000, 1024, NewXCondBlockStrategy()},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.loop; i++ {
				_, _ = reportFile.WriteString(fmt.Sprintf("\n====== begin string report(%s, %d) ======\n", tc.name, i))
				testXSinglePipelineDisruptorString(t, tc.gTotal, tc.tasks, tc.capacity, tc.bs, true, reportFile, errorCounter)
			}
		})
	}
	t.Logf("errors: %d\n", errorCounter.Load())
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
		name   string
		gTotal int
		tasks  int
	}{
		{"nochan 10*100", 10, 100},
		{"nochan 100*10000", 100, 10000},
		{"nochan 500*10000", 500, 10000},
		{"nochan 1000*10000", 1000, 10000},
		{"nochan 5000*10000", 5000, 10000},
		{"nochan 10000*10000", 10000, 10000},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			testNoCacheChannel(t, 0, tc.gTotal, tc.tasks)
		})
	}
}

func TestCacheChannel(t *testing.T) {
	testcases := []struct {
		name   string
		gTotal int
		tasks  int
	}{
		{"cachechan 10*100", 10, 100},
		{"cachechan 100*10000", 100, 10000},
		{"cachechan 500*10000", 500, 10000},
		{"cachechan 1000*10000", 1000, 10000},
		{"cachechan 5000*10000", 5000, 10000},
		{"cachechan 10000*10000", 10000, 10000},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			testNoCacheChannel(t, 1024*1024, tc.gTotal, tc.tasks)
		})
	}
}

func testXSinglePipelineDisruptorWithRandomSleep(t *testing.T, num, capacity int) {
	wg := &sync.WaitGroup{}
	wg.Add(num)
	results := map[string]struct{}{}
	disruptor := NewXSinglePipelineDisruptor[string](uint64(capacity),
		NewXCacheChannelBlockStrategy(),
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
