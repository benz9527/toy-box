package queue

type RingBufferEntry[T any] interface {
	GetValue() T
	GetCursor() uint64
	Store(cursor uint64, value T)
}

type RingBuffer[T any] interface {
	Capacity() uint64
	LoadEntry(cursor uint64) (RingBufferEntry[T], bool)
	StoreEntry(cursor uint64, value T)
}

type RingBufferCursor interface {
	Next() uint64
	NextN(n uint64) uint64
	Load() uint64
	CAS(old, new uint64) bool
}
