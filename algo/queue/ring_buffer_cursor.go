package queue

import "sync/atomic"

var (
	_ RingBufferCursor = (*xRingBufferCursor)(nil)
)

// xRingBufferCursor is a cursor for xRingBuffer.
// Only increase, if it overflows, it will be reset to 0.
// Occupy a whole cache line (flag+tag+data) and a cache line data is 64 bytes.
// L1D cache: cat /sys/devices/system/cpu/cpu0/cache/index0/coherency_line_size
// L1I cache: cat /sys/devices/system/cpu/cpu0/cache/index1/coherency_line_size
// L2 cache: cat /sys/devices/system/cpu/cpu0/cache/index2/coherency_line_size
// L3 cache: cat /sys/devices/system/cpu/cpu0/cache/index3/coherency_line_size
// MESI (Modified-Exclusive-Shared-Invalid)
// RAM data -> L3 cache -> L2 cache -> L1 cache -> CPU register
// CPU register (cache hit) -> L1 cache -> L2 cache -> L3 cache -> RAM data
type xRingBufferCursor struct {
	// sequence consistency data race free program
	// avoid load into cpu cache will be broken by others data
	// to compose a data race cache line
	_      [7]uint64 // padding for CPU cache line, avoid false sharing
	cursor uint64    // space waste to exchange for performance
	_      [7]uint64 // padding for CPU cache line, avoid false sharing
}

func NewXRingBufferCursor() RingBufferCursor {
	return &xRingBufferCursor{}
}

func (x *xRingBufferCursor) Increase() uint64 {
	return atomic.AddUint64(&x.cursor, 1)
}

func (x *xRingBufferCursor) AtomicLoad() uint64 {
	return atomic.LoadUint64(&x.cursor)
}

func (x *xRingBufferCursor) CASStore(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&x.cursor, old, new)
}
