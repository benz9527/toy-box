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

func (pub *xSinglePipelinePublisher[T]) Start() error {
	if atomic.CompareAndSwapInt32((*int32)(&pub.status), int32(pubReady), int32(pubRunning)) {
		return nil
	}
	return fmt.Errorf("publisher already started")
}

func (pub *xSinglePipelinePublisher[T]) Stop() error {
	if atomic.CompareAndSwapInt32((*int32)(&pub.status), int32(pubRunning), int32(pubReady)) {
		return nil
	}
	return fmt.Errorf("publisher already stopped")
}

func (pub *xSinglePipelinePublisher[T]) IsStopped() bool {
	return atomic.LoadInt32((*int32)(&pub.status)) == int32(pubReady)
}

func (pub *xSinglePipelinePublisher[T]) Publish(event T) (uint64, bool, error) {
	if pub.IsStopped() {
		return 0, false, fmt.Errorf("publisher closed")
	}
	nextWriteCursor := pub.seq.GetWriteCursor().Increase()
	for {
		readCursor := pub.seq.LoadReadCursor() - 1
		if nextWriteCursor <= readCursor+pub.seq.Capacity() {
			pub.rb.StoreElement(nextWriteCursor-1, event)
			pub.strategy.Done()
			return nextWriteCursor, true, nil
		}
		runtime.Gosched()
		if pub.IsStopped() {
			return 0, false, fmt.Errorf("publisher closed")
		}
	}
}

func (pub *xSinglePipelinePublisher[T]) PublishTimeout(event T, timeout time.Duration) (uint64, bool, error) {
	if pub.IsStopped() {
		return 0, false, fmt.Errorf("publisher closed")
	}
	nextCursor := pub.seq.GetWriteCursor().Increase()
	ok := pub.publishAt(event, nextCursor)
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
			if ok = pub.publishAt(event, nextCursor); ok {
				return nextCursor, true, nil
			}
			runtime.Gosched()
		}
		if pub.IsStopped() {
			return 0, false, fmt.Errorf("publisher closed")
		}
	}
}

// unstable result under concurrent scenario
func (pub *xSinglePipelinePublisher[T]) tryWriteWindow() int {
	tryNext := pub.seq.GetWriteCursor().AtomicLoad() + 1
	readCursor := pub.seq.LoadReadCursor()
	if tryNext < readCursor+pub.capacity {
		return int(readCursor + pub.seq.Capacity() - tryNext)
	}
	return -int(tryNext - (readCursor + pub.seq.Capacity()))
}

func (pub *xSinglePipelinePublisher[T]) publishAt(event T, cursor uint64) bool {
	idx := pub.seq.LoadReadCursor() - 1
	if cursor > idx+pub.seq.Capacity() {
		return false
	}
	pub.rb.StoreElement(cursor-1, event)
	pub.strategy.Done()
	return true
}
