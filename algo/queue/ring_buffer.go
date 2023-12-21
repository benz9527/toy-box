//go:build linux

package queue

import (
	"sync/atomic"
)

var (
	_ RingBufferElement[struct{}] = (*xRingBufferElement[struct{}])(nil)
	_ RingBuffer[struct{}]        = (*xRingBuffer[struct{}])(nil)
)

type xRingBufferElement[T any] struct {
	// index of the element in the ring buffer and
	// it is also the "lock" to protect the value
	// in lock-free mode by atomic operation
	cursor uint64
	value  T
}

func (e *xRingBufferElement[T]) GetValue() T {
	return e.value
}

func (e *xRingBufferElement[T]) GetCursor() uint64 {
	return atomic.LoadUint64(&e.cursor)
}

type xRingBuffer[T any] struct {
	capacityMask uint64
	buffer       []*xRingBufferElement[T]
	valueGuard   T
}

func NewXRingBuffer[T any](capacity uint64) RingBuffer[T] {
	rb := &xRingBuffer[T]{
		capacityMask: capacity - 1,
		buffer:       make([]*xRingBufferElement[T], capacity),
		valueGuard:   *new(T),
	}
	for i := uint64(0); i < capacity; i++ {
		rb.buffer[i] = &xRingBufferElement[T]{}
	}
	return rb
}

func (rb *xRingBuffer[T]) Capacity() uint64 {
	return rb.capacityMask + 1
}

func (rb *xRingBuffer[T]) StoreElement(cursor uint64, value T) {
	e := rb.buffer[cursor&rb.capacityMask] // get remainder
	// atomic operation should be called at the end of the function
	// otherwise, the value of cursor may be changed by other goroutines
	e.value = value
	atomic.StoreUint64(&e.cursor, cursor)
}

func (rb *xRingBuffer[T]) LoadElement(cursor uint64) (RingBufferElement[T], bool) {
	e := rb.buffer[cursor&rb.capacityMask]
	if e.GetCursor() == cursor {
		return e, true
	}
	return &xRingBufferElement[T]{
		cursor: e.GetCursor(),
		value:  rb.valueGuard,
	}, false
}
