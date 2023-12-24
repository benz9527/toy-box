package queue

type RingBufferElement[T any] interface {
	GetValue() T
	GetCursor() uint64
}

type RingBuffer[T any] interface {
	Capacity() uint64
	LoadElement(cursor uint64) (RingBufferElement[T], bool)
	StoreElement(cursor uint64, value T)
	Free()
}

type RingBufferCursor interface {
	Increase() uint64
	Load() uint64
	CAS(old, new uint64) bool
}
