//go:build linux && !386

package pubsub

import (
	"fmt"
	"github.com/benz9527/toy-box/algo/bit"
	"github.com/benz9527/toy-box/algo/queue"
	"sync/atomic"
	"time"
)

type disruptorStatus int32

const (
	disruptorReady disruptorStatus = iota
	disruptorRunning
)

type xSinglePipelineDisruptor[T any] struct {
	pub interface {
		Publisher[T]
		stopper
	}
	sub interface {
		Subscriber[T]
		stopper
	}
	status disruptorStatus
}

func NewXSinglePipelineDisruptor[T any](
	capacity uint64,
	strategy BlockStrategy,
	handler EventHandler[T],
) Disruptor[T] {
	capacity = bit.RoundupPowOf2ByCeil(capacity)
	seq := NewXSequencer(capacity)
	// Can't start from 0, because 0 will be treated as nil value
	seq.GetWriteCursor().Next()
	seq.GetReadCursor().Next()
	rb := queue.NewXRingBuffer[T](capacity)
	pub := newXSinglePipelinePublisher[T](seq, rb, strategy)
	sub := newXSinglePipelineSubscriber[T](rb, handler, seq, strategy)
	d := &xSinglePipelineDisruptor[T]{
		pub:    pub,
		sub:    sub,
		status: disruptorReady,
	}
	return d
}

func (dis *xSinglePipelineDisruptor[T]) Start() error {
	if atomic.CompareAndSwapInt32((*int32)(&dis.status), int32(disruptorReady), int32(disruptorRunning)) {
		if err := dis.sub.Start(); err != nil {
			atomic.StoreInt32((*int32)(&dis.status), int32(disruptorReady))
			return err
		}
		if err := dis.pub.Start(); err != nil {
			atomic.StoreInt32((*int32)(&dis.status), int32(disruptorReady))
			return err
		}
		return nil
	}
	return fmt.Errorf("disruptor already started")
}

func (dis *xSinglePipelineDisruptor[T]) Stop() error {
	if atomic.CompareAndSwapInt32((*int32)(&dis.status), int32(disruptorRunning), int32(disruptorReady)) {
		if err := dis.pub.Stop(); err != nil {
			atomic.CompareAndSwapInt32((*int32)(&dis.status), int32(disruptorRunning), int32(disruptorReady))
			return err
		}
		if err := dis.sub.Stop(); err != nil {
			atomic.CompareAndSwapInt32((*int32)(&dis.status), int32(disruptorRunning), int32(disruptorReady))
			return err
		}
		return nil
	}
	return fmt.Errorf("disruptor already stopped")
}

func (dis *xSinglePipelineDisruptor[T]) IsStopped() bool {
	return atomic.LoadInt32((*int32)(&dis.status)) != int32(disruptorRunning)
}

func (dis *xSinglePipelineDisruptor[T]) Publish(event T) (uint64, bool, error) {
	return dis.pub.Publish(event)
}

func (dis *xSinglePipelineDisruptor[T]) PublishTimeout(event T, timeout time.Duration) {
	dis.pub.PublishTimeout(event, timeout)
}

func (dis *xSinglePipelineDisruptor[T]) RegisterSubscriber(sub Subscriber[T]) error {
	// Single pipeline disruptor only support one subscriber to consume the events.
	// It will be registered at the construction.
	return nil
}
