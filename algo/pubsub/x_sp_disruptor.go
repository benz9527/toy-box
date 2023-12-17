//go:build linux && !386

package pubsub

import (
	"fmt"
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
	capacity = ceilCapacity(capacity)
	seq := NewXSequencer(capacity)
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

func (x *xSinglePipelineDisruptor[T]) Start() error {
	if atomic.CompareAndSwapInt32((*int32)(&x.status), int32(disruptorReady), int32(disruptorRunning)) {
		if err := x.sub.Start(); err != nil {
			atomic.StoreInt32((*int32)(&x.status), int32(disruptorReady))
			return err
		}
		if err := x.pub.Start(); err != nil {
			atomic.StoreInt32((*int32)(&x.status), int32(disruptorReady))
			return err
		}
		return nil
	}
	return fmt.Errorf("disruptor already started")
}

func (x *xSinglePipelineDisruptor[T]) Stop() error {
	if atomic.CompareAndSwapInt32((*int32)(&x.status), int32(disruptorRunning), int32(disruptorReady)) {
		if err := x.pub.Stop(); err != nil {
			atomic.CompareAndSwapInt32((*int32)(&x.status), int32(disruptorRunning), int32(disruptorReady))
			return err
		}
		if err := x.sub.Stop(); err != nil {
			atomic.CompareAndSwapInt32((*int32)(&x.status), int32(disruptorRunning), int32(disruptorReady))
			return err
		}
		return nil
	}
	return fmt.Errorf("disruptor already stopped")
}

func (x *xSinglePipelineDisruptor[T]) IsStopped() bool {
	return atomic.LoadInt32((*int32)(&x.status)) != int32(disruptorRunning)
}

func (x *xSinglePipelineDisruptor[T]) Publish(event T) (uint64, bool, error) {
	return x.pub.Publish(event)
}

func (x *xSinglePipelineDisruptor[T]) PublishTimeout(event T, timeout time.Duration) (uint64, bool, error) {
	return x.pub.PublishTimeout(event, timeout)
}

func (x *xSinglePipelineDisruptor[T]) RegisterSubscriber(sub Subscriber[T]) error {
	// Single pipeline disruptor only support one subscriber to consume the events.
	// It will be registered at the construction.
	return nil
}

func ceilCapacity(capacity uint64) uint64 {
	if capacity&(capacity-1) == 0 {
		return capacity
	}
	var _cap uint64 = 1
	for _cap < capacity {
		_cap <<= 1
	}
	return _cap
}
