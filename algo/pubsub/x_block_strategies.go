package pubsub

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	_ "unsafe"
)

var (
	_ BlockStrategy = (*xGoSchedBlockStrategy)(nil)
	_ BlockStrategy = (*xSleepBlockStrategy)(nil)
	_ BlockStrategy = (*xCpuNoOpLoopBlockStrategy)(nil)
	_ BlockStrategy = (*xOsYieldBlockStrategy)(nil)
	_ BlockStrategy = (*xNoCacheChannelBlockStrategy)(nil)
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

func (x *xSleepBlockStrategy) WaitFor(eqFn func() bool) {
	time.Sleep(x.sleepTime)
}

func (x *xSleepBlockStrategy) Done() {
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

func (x *xCpuNoOpLoopBlockStrategy) WaitFor(eqFn func() bool) {
	procYield(x.cycles)
}

func (x *xCpuNoOpLoopBlockStrategy) Done() {}

//go:linkname osYield runtime.osyield
func osYield()

type xOsYieldBlockStrategy struct{}

func NewXOsYieldBlockStrategy() BlockStrategy {
	return &xOsYieldBlockStrategy{}
}

func (x *xOsYieldBlockStrategy) WaitFor(fn func() bool) {
	osYield()
}

func (x *xOsYieldBlockStrategy) Done() {}

type xNoCacheChannelBlockStrategy struct {
	ch     chan struct{}
	status *atomic.Bool
}

func NewXNoCacheChannelBlockStrategy() BlockStrategy {
	return &xNoCacheChannelBlockStrategy{
		ch:     make(chan struct{}),
		status: &atomic.Bool{},
	}
}

func (x *xNoCacheChannelBlockStrategy) WaitFor(eqFn func() bool) {
	// Try to block
	if x.status.CompareAndSwap(false, true) {
		// Double check
		if !eqFn() {
			// Block, wait for signal
			<-x.ch
		} else {
			//  Double check failed, reset status
			if !x.status.CompareAndSwap(true, false) {
				// Wait for release
				<-x.ch
			}
			return
		}
	}
}

func (x *xNoCacheChannelBlockStrategy) Done() {
	// Release
	if x.status.CompareAndSwap(true, false) {
		// Send signal
		x.ch <- struct{}{}
	}
	return
}

type xCondBlockStrategy struct {
	cond *sync.Cond
}

func NewXCondBlockStrategy() BlockStrategy {
	return &xCondBlockStrategy{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (x *xCondBlockStrategy) WaitFor(eqFn func() bool) {
	x.cond.L.Lock()
	defer x.cond.L.Unlock()
	if eqFn() {
		return
	}
	x.cond.Wait()
}

func (x *xCondBlockStrategy) Done() {
	x.cond.L.Lock()
	defer x.cond.L.Unlock()
	x.cond.Broadcast()
}
