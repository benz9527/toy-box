package container

import "io"

type IQueueItem interface {
	ObjectHash() string
}

type DequeueProcessFunc[T IQueueItem, V any] func(item T) (any, error)

// DequeueProcessor is a queue with dequeue process function.
// Note: It seems useless, but maybe useful for some cases.
type DequeueProcessor[T IQueueItem] interface {
	DequeueProcess(fn DequeueProcessFunc[T, any]) (any, error)
}

type IQueue[T IQueueItem] interface {
	// Closer Closing and cleaning the queue
	io.Closer
	// Enqueue an item into the queue (tail).
	Enqueue(item T) error
	// Dequeue an item from the queue (front).
	Dequeue() (T, error)
}

// IZQueue is a queue with zero duplication of queue item.
type IZQueue[T IQueueItem] interface {
	IQueue[T]
	// EnqueueIfNotPresent Checking item by ObjectHash() with
	// specific hash function and add it if not present
	EnqueueIfNotPresent(item T) error
}

// IDeque is a double-ended queue.
type IDeque[T IQueueItem] interface {
	IQueue[T]
	// EnqueueFront an item into the queue (front).
	EnqueueFront(item T) error
	// EnqueueTail an item into the queue (tail).
	EnqueueTail(item T) error
	// DequeueFront an item from the queue (front).
	DequeueFront() (T, error)
	// DequeueTail an item from the queue (tail).
	DequeueTail() (T, error)
}

// IZDeque is a double-ended queue with zero duplication of queue item.
type IZDeque[T IQueueItem] interface {
	IDeque[T]
	// EnqueueFrontIfNotPresent Checking item by ObjectHash() with
	// specific hash function and add it if not present
	EnqueueFrontIfNotPresent(item T) error
}
