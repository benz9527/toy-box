package pubsub

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/benz9527/toy-box/algo/queue"
)

type publisherStatus int32

const (
	pubReady publisherStatus = iota
	pubRunning
)

var (
	_ Publisher[int] = (*xSinglePipelinePublisher[int])(nil)
	_ stopper        = (*xSinglePipelinePublisher[int])(nil)
)

type xSinglePipelinePublisher[T any] struct {
	seq      Sequencer
	rb       queue.RingBuffer[T]
	strategy BlockStrategy
	capacity uint64
	status   publisherStatus
}

func newXSinglePipelinePublisher[T any](seq Sequencer, rb queue.RingBuffer[T], strategy BlockStrategy) *xSinglePipelinePublisher[T] {
	return &xSinglePipelinePublisher[T]{
		seq:      seq,
		rb:       rb,
		strategy: strategy,
		capacity: seq.Capacity(),
		status:   pubReady,
	}
}

func (x *xSinglePipelinePublisher[T]) Start() error {
	if atomic.CompareAndSwapInt32((*int32)(&x.status), int32(pubReady), int32(pubRunning)) {
		return nil
	}
	return fmt.Errorf("publisher already started")
}

func (x *xSinglePipelinePublisher[T]) Stop() error {
	if atomic.CompareAndSwapInt32((*int32)(&x.status), int32(pubRunning), int32(pubReady)) {
		return nil
	}
	return fmt.Errorf("publisher already stopped")
}

func (x *xSinglePipelinePublisher[T]) IsStopped() bool {
	return atomic.LoadInt32((*int32)(&x.status)) == int32(pubReady)
}

func (x *xSinglePipelinePublisher[T]) Publish(event T) (uint64, bool, error) {
	if x.IsStopped() {
		return 0, false, fmt.Errorf("publisher closed")
	}
	nextCursor := x.seq.GetWriteCursor().Increase()
	for {
		readCursor := x.seq.LoadReadCursor() - 1
		if nextCursor <= readCursor+x.seq.Capacity() {
			x.rb.StoreElement(nextCursor-1, event)
			x.strategy.Done()
			return nextCursor, true, nil
		}
		runtime.Gosched()
		if x.IsStopped() {
			return 0, false, fmt.Errorf("publisher closed")
		}
	}

}

func (x *xSinglePipelinePublisher[T]) PublishTimeout(event T, timeout time.Duration) (uint64, bool, error) {
	if x.IsStopped() {
		return 0, false, fmt.Errorf("publisher closed")
	}
	nextCursor := x.seq.GetWriteCursor().Increase()
	ok := x.publishAt(event, nextCursor)
	if ok {
		return nextCursor, true, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return 0, false, fmt.Errorf("publish timeout")
		default:
			if ok = x.publishAt(event, nextCursor); ok {
				return nextCursor, true, nil
			}
			runtime.Gosched()
		}
		if x.IsStopped() {
			return 0, false, fmt.Errorf("publisher closed")
		}
	}
}

// unstable result under concurrent scenario
func (x *xSinglePipelinePublisher[T]) tryWriteWindow() int {
	tryNext := x.seq.GetWriteCursor().AtomicLoad() + 1
	readCursor := x.seq.LoadReadCursor()
	if tryNext < readCursor+x.capacity {
		return int(readCursor + x.seq.Capacity() - tryNext)
	}
	return -int(tryNext - (readCursor + x.seq.Capacity()))
}

func (x *xSinglePipelinePublisher[T]) publishAt(event T, cursor uint64) bool {
	idx := x.seq.LoadReadCursor() - 1
	if cursor > idx+x.seq.Capacity() {
		return false
	}
	x.rb.StoreElement(cursor-1, event)
	x.strategy.Done()
	return true
}
