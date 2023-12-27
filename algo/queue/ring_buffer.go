//go:build linux

package queue

// circular queue

import (
	"github.com/benz9527/toy-box/algo/bit"
	"runtime"
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
	buffer       []*xRingBufferElement[T] // []atomic.Pointer[xRingBufferElement[T]]
	valueGuard   T
}

const _100MB = 100 * 1024 * 1024

func NewXRingBuffer[T any](capacity uint64) RingBuffer[T] {
	if capacity > _100MB {
		panic("capacity is too large")
	}
	if bit.IsPowOf2(capacity) {
		capacity = bit.RoundupPowOf2ByCeil(capacity)
		if capacity > _100MB {
			panic("capacity is too large")
		}
	}

	rb := &xRingBuffer[T]{
		capacityMask: capacity - 1,
		buffer:       make([]*xRingBufferElement[T], capacity),
		valueGuard:   *new(T),
	}
	for i := uint64(0); i < capacity; i++ {
		rb.buffer[i] = &xRingBufferElement[T]{}
	}
	runtime.SetFinalizer(rb, func(rb *xRingBuffer[T]) {
		clear(rb.buffer)
	})
	return rb
}

func (rb *xRingBuffer[T]) Capacity() uint64 {
	return rb.capacityMask + 1
}

func (rb *xRingBuffer[T]) StoreElement(cursor uint64, value T) {
	e := rb.buffer[cursor&rb.capacityMask]
	// atomic operation should be called at the end of the function
	// otherwise, the value of cursor may be changed by other goroutines
	// Go1.19 atomic guarantees the sequentially consistent order
	// Go1.18 atomic.Pointer[T]
	e.value = value
	atomic.StoreUint64(&e.cursor, cursor)
}

func (rb *xRingBuffer[T]) LoadElement(cursor uint64) (RingBufferElement[T], bool) {
	e := rb.buffer[cursor&rb.capacityMask]
	if e.GetCursor() == cursor {
		return e, true
	}
	return &xRingBufferElement[T]{
		cursor: 0,
		value:  rb.valueGuard,
	}, false
}
