package pubsub

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/benz9527/toy-box/algo/queue"
)

type subscriberStatus int32

const (
	subReady subscriberStatus = iota
	subRunning
)

const (
	activeSpin  = 4
	passiveSpin = 2
)

type xSinglePipelineSubscriber[T any] struct {
	status   subscriberStatus
	seq      Sequencer
	rb       queue.RingBuffer[T]
	strategy BlockStrategy
	handler  EventHandler[T]
	spin     int32
}

func newXSinglePipelineSubscriber[T any](
	rb queue.RingBuffer[T],
	handler EventHandler[T],
	seq Sequencer,
	strategy BlockStrategy,
) *xSinglePipelineSubscriber[T] {
	ncpu := runtime.NumCPU()
	spin := 0
	if ncpu > 1 {
		spin = activeSpin
	}
	return &xSinglePipelineSubscriber[T]{
		status:   subReady,
		seq:      seq,
		rb:       rb,
		strategy: strategy,
		handler:  handler,
		spin:     int32(spin),
	}
}

func (x *xSinglePipelineSubscriber[T]) Start() error {
	if atomic.CompareAndSwapInt32((*int32)(&x.status), int32(subReady), int32(subRunning)) {
		go x.eventsHandle()
		return nil
	}
	return fmt.Errorf("subscriber already started")
}

func (x *xSinglePipelineSubscriber[T]) Stop() error {
	if atomic.CompareAndSwapInt32((*int32)(&x.status), int32(subRunning), int32(subReady)) {
		x.strategy.Done()
		return nil
	}
	return fmt.Errorf("subscriber already stopped")
}

func (x *xSinglePipelineSubscriber[T]) IsStopped() bool {
	return atomic.LoadInt32((*int32)(&x.status)) == int32(subReady)
}

func (x *xSinglePipelineSubscriber[T]) eventsHandle() {
	readCursor := x.seq.LoadReadCursor()
	for {
		if x.IsStopped() {
			return
		}
		i := int32(0)
		for {
			if x.IsStopped() {
				return
			}
			if e, exists := x.rb.LoadElement(readCursor - 1); exists {
				readCursor = x.seq.NextReadCursor()
				// FIXME handle error
				_ = x.HandleEvent(e.GetValue())
				i = 0
				break
			} else {
				if i < atomic.LoadInt32(&x.spin) {
					procYield(30)
				} else if i < atomic.LoadInt32(&x.spin)+passiveSpin {
					runtime.Gosched()
				} else {
					x.strategy.WaitFor(func() bool {
						return e.GetCursor() == readCursor
					})
					i = 0
				}
				i++
			}
		}
	}
}

func (x *xSinglePipelineSubscriber[T]) HandleEvent(event T) error {
	defer x.strategy.Done()
	err := x.handler(event)
	return err
}
