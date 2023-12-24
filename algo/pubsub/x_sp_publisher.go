package pubsub

import (
	"context"
	"fmt"
	"log/slog"
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
		readCursor := pub.seq.GetReadCursor().Load()
		if nextWriteCursor-readCursor <= pub.capacity {
			pub.rb.StoreElement(nextWriteCursor-1, event)
			pub.strategy.Done()
			return nextWriteCursor - 1, true, nil
		} else {
			pub.strategy.Done()
		}
		runtime.Gosched()
		if pub.IsStopped() {
			return 0, false, fmt.Errorf("publisher closed")
		}
	}
}

func (pub *xSinglePipelinePublisher[T]) PublishTimeout(event T, timeout time.Duration) {
	go func() {
		if pub.IsStopped() {
			slog.Warn("publisher closed", "event", event)
			return
		}
		nextCursor := pub.seq.GetWriteCursor().Increase()
		var ok bool
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				slog.Warn("publish timeout", "event", event)
				return
			default:
				if ok = pub.publishAt(event, nextCursor-1); ok {
					slog.Info("publish success", "event", event)
					return
				}
			}
			runtime.Gosched()
			if pub.IsStopped() {
				slog.Warn("publisher closed", "event", event)
				return
			}
		}
	}()
}

// unstable result under concurrent scenario
func (pub *xSinglePipelinePublisher[T]) tryWriteWindow() int {
	tryNext := pub.seq.GetWriteCursor().Load() + 1
	readCursor := pub.seq.GetReadCursor().Load() - 1
	if tryNext < readCursor+pub.capacity {
		return int(readCursor + pub.capacity - tryNext)
	}
	return -int(tryNext - (readCursor + pub.capacity))
}

func (pub *xSinglePipelinePublisher[T]) publishAt(event T, cursor uint64) bool {
	readCursor := pub.seq.GetReadCursor().Load() - 1
	if cursor > readCursor+pub.seq.Capacity() {
		return false
	}
	pub.rb.StoreElement(cursor, event)
	pub.strategy.Done()
	return true
}
