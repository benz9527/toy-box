package concurrency

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type GenericWaitChannel[T comparable] interface {
	Wait() <-chan T
}

type waitChannel[T comparable] struct {
	ch <-chan T
}

func (w *waitChannel[T]) Wait() <-chan T {
	return w.ch
}

func WaitChannelWrapper[T comparable](ch <-chan T) GenericWaitChannel[T] {
	return &waitChannel[T]{ch: ch}
}

var (
	neverStopWaitC = WaitChannelWrapper[struct{}](make(chan struct{}))
)

type GenericChannel[T comparable] interface {
	io.Closer
	GenericWaitChannel[T]
	IsClosed() bool
	Send(v T, nonBlocking ...bool) error
}

// safeChannel is a generic channel wrapper.
// Why we need this wrapper? For the following reasons:
// 1. We need to make sure the channel is closed only once.
type safeChannel[T comparable] struct {
	queueC chan T // Receive data to temporary queue.
	// cachedOneC chan T // Delivery data from temporary queue.
	isClosed  atomic.Bool
	isClosing atomic.Bool
	once      *sync.Once
}

var (
	_ GenericWaitChannel[struct{}] = &safeChannel[struct{}]{} // type check assertion
)

func NewSafeChannel[T comparable](chSize ...int) GenericChannel[T] {
	isNoCacheCh := true
	size := 1
	if len(chSize) > 0 {
		if chSize[0] > 0 {
			size = chSize[0]
			isNoCacheCh = false
		}
	}
	if isNoCacheCh {
		return &safeChannel[T]{
			queueC: make(chan T),
			once:   &sync.Once{},
		}
	}
	return &safeChannel[T]{
		queueC: make(chan T, size),
		once:   &sync.Once{},
	}
}

func (c *safeChannel[T]) IsClosed() bool {
	return c.isClosed.Load()
}

// Close According to the Go memory model, a send on a channel happens before
// the corresponding receive from that channel completes
// https://go.dev/doc/articles/race_detector
func (c *safeChannel[T]) Close() error {
	c.once.Do(func() {
		c.isClosing.Store(true)
		// FIXME data race, it may do harm to the application
		close(c.queueC)
		c.isClosed.Store(true)
		c.isClosing.Store(false)
	})
	return nil
}

func (c *safeChannel[T]) Wait() <-chan T {
	return c.queueC
}

func (c *safeChannel[T]) Send(v T, nonBlocking ...bool) error {
	if c.isClosing.Load() {
		return fmt.Errorf("channel is closing")
	}
	if c.isClosed.Load() {
		return fmt.Errorf("channel has been closed")
	}

	if len(nonBlocking) <= 0 {
		nonBlocking = []bool{false}
	}
	if !nonBlocking[0] {
		c.queueC <- v
	} else {
		// non blocking send
		select {
		case c.queueC <- v:
		default:

		}
	}
	return nil
}

func ContextForChannel(parentC GenericWaitChannel[struct{}]) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-parentC.Wait():
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}
