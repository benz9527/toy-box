package queue

type RingBufferElement[T any] interface {
	GetValue() T
	setValue(v T)
	GetCursor() uint64
	setCursor(cursor uint64)
}

type RingBuffer[T any] interface {
	Capacity() uint64
	LoadElement(cursor uint64) (RingBufferElement[T], bool)
	StoreElement(cursor uint64, value T)
}

type RingBufferCursor interface {
	Increase() uint64
	AtomicLoad() uint64
	CASStore(old, new uint64) bool
}
