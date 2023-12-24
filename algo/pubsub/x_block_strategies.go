package pubsub

import (
	"github.com/benz9527/toy-box/algo/queue"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
	_ "unsafe"
)

var (
	_ BlockStrategy = (*xGoSchedBlockStrategy)(nil)
	_ BlockStrategy = (*xSleepBlockStrategy)(nil)
	_ BlockStrategy = (*xCpuNoOpLoopBlockStrategy)(nil)
	_ BlockStrategy = (*xOsYieldBlockStrategy)(nil)
	_ BlockStrategy = (*xCacheChannelBlockStrategy)(nil)
	_ BlockStrategy = (*xCondBlockStrategy)(nil)
)

type xGoSchedBlockStrategy struct{}

func NewXGoSchedBlockStrategy() BlockStrategy {
	return &xGoSchedBlockStrategy{}
}

func (x *xGoSchedBlockStrategy) WaitFor(eqFn func() bool) {
	// Let go runtime schedule other goroutines
	// Current goroutine will yield CPU
	runtime.Gosched()
}

func (x *xGoSchedBlockStrategy) Done() {
	// do nothing
}

type xSleepBlockStrategy struct {
	sleepTime time.Duration
}

func NewXSleepBlockStrategy(sleepTime time.Duration) BlockStrategy {
	// around 10us, cpu load 2-3%
	// eq to 5us, cpu load 10%
	// lt 5us, cpu load 100%
	return &xSleepBlockStrategy{
		sleepTime: sleepTime,
	}
}

func (bs *xSleepBlockStrategy) WaitFor(eqFn func() bool) {
	time.Sleep(bs.sleepTime)
}

func (bs *xSleepBlockStrategy) Done() {
	// do nothing
}

type xCpuNoOpLoopBlockStrategy struct {
	cycles uint32
}

//go:linkname procYield runtime.procyield
func procYield(cycles uint32)

func NewXCpuNoOpLoopBlockStrategy(cycles uint32) BlockStrategy {
	return &xCpuNoOpLoopBlockStrategy{
		cycles: cycles,
	}
}

func (bs *xCpuNoOpLoopBlockStrategy) WaitFor(eqFn func() bool) {
	procYield(bs.cycles)
}

func (bs *xCpuNoOpLoopBlockStrategy) Done() {}

//go:linkname osYield runtime.osyield
func osYield()

type xOsYieldBlockStrategy struct{}

func NewXOsYieldBlockStrategy() BlockStrategy {
	return &xOsYieldBlockStrategy{}
}

func (bs *xOsYieldBlockStrategy) WaitFor(fn func() bool) {
	osYield()
}

func (bs *xOsYieldBlockStrategy) Done() {}

type xCacheChannelBlockStrategy struct {
	_      [queue.CacheLinePadSize - unsafe.Sizeof(*new(uint64))]byte
	status uint64
	_      [queue.CacheLinePadSize - unsafe.Sizeof(*new(uint64))]byte
	ch     chan struct{}
}

func NewXCacheChannelBlockStrategy() BlockStrategy {
	return &xCacheChannelBlockStrategy{
		ch:     make(chan struct{}, 1),
		status: 0,
	}
}

func (bs *xCacheChannelBlockStrategy) WaitFor(eqFn func() bool) {
	// Try to block
	if atomic.CompareAndSwapUint64(&bs.status, 0, 1) {
		// Double check
		if !eqFn() {
			// Block, wait for signal
			<-bs.ch
		} else {
			//  Double check failed, reset status
			if !atomic.CompareAndSwapUint64(&bs.status, 1, 0) {
				// Wait for release
				<-bs.ch
			}
		}
	}
}

func (bs *xCacheChannelBlockStrategy) Done() {
	// Release
	if atomic.CompareAndSwapUint64(&bs.status, 1, 0) {
		// Send signal
		bs.ch <- struct{}{}
	}
}

type xCondBlockStrategy struct {
	cond *sync.Cond
}

func NewXCondBlockStrategy() BlockStrategy {
	return &xCondBlockStrategy{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (bs *xCondBlockStrategy) WaitFor(eqFn func() bool) {
	bs.cond.L.Lock()
	defer bs.cond.L.Unlock()
	if eqFn() {
		return
	}
	bs.cond.Wait()
}

func (bs *xCondBlockStrategy) Done() {
	bs.cond.Broadcast()
}
