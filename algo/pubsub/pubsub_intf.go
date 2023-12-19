package pubsub

import (
	"time"

	"github.com/benz9527/toy-box/algo/queue"
)

type stopper interface {
	Start() error
	Stop() error
	IsStopped() bool
}

type Publisher[T any] interface {
	Publish(event T) (uint64, bool, error)
	PublishTimeout(event T, timeout time.Duration)
}

type Producer[T any] Publisher[T]

type BlockStrategy interface {
	WaitFor(eqFn func() bool)
	Done()
}

type EventHandler[T any] func(event T) error // OnEvent

type Subscriber[T any] interface {
	HandleEvent(event T) error
}

type Sequencer interface {
	NextReadCursor() uint64
	LoadReadCursor() uint64
	Capacity() uint64
	GetWriteCursor() queue.RingBufferCursor
}

type Disruptor[T any] interface {
	Publisher[T]
	stopper
	RegisterSubscriber(sub Subscriber[T]) error
}
