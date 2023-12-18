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
	rb       queue.RingBuffer[T]
	seq      Sequencer
	strategy BlockStrategy
	handler  EventHandler[T]
	status   subscriberStatus
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

func (sub *xSinglePipelineSubscriber[T]) Start() error {
	if atomic.CompareAndSwapInt32((*int32)(&sub.status), int32(subReady), int32(subRunning)) {
		go sub.eventsHandle()
		return nil
	}
	return fmt.Errorf("subscriber already started")
}

func (sub *xSinglePipelineSubscriber[T]) Stop() error {
	if atomic.CompareAndSwapInt32((*int32)(&sub.status), int32(subRunning), int32(subReady)) {
		sub.strategy.Done()
		return nil
	}
	return fmt.Errorf("subscriber already stopped")
}

func (sub *xSinglePipelineSubscriber[T]) IsStopped() bool {
	return atomic.LoadInt32((*int32)(&sub.status)) == int32(subReady)
}

func (sub *xSinglePipelineSubscriber[T]) eventsHandle() {
	readCursor := sub.seq.LoadReadCursor()
	for {
		if sub.IsStopped() {
			return
		}
		i := int32(0)
		for {
			if sub.IsStopped() {
				return
			}
			if e, exists := sub.rb.LoadElement(readCursor - 1); exists {
				readCursor = sub.seq.NextReadCursor()
				// FIXME handle error
				_ = sub.HandleEvent(e.GetValue())
				i = 0
				break
			} else {
				if i < sub.spin {
					procYield(30)
				} else if i < sub.spin+passiveSpin {
					runtime.Gosched()
				} else {
					sub.strategy.WaitFor(func() bool {
						return e.GetCursor() == readCursor
					})
					i = 0
				}
				i++
			}
		}
	}
}

func (sub *xSinglePipelineSubscriber[T]) HandleEvent(event T) error {
	defer sub.strategy.Done()
	err := sub.handler(event)
	return err
}
