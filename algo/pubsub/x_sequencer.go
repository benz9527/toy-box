package pubsub

import (
	"sync/atomic"

	"github.com/benz9527/toy-box/algo/queue"
)

type xSequencer struct {
	writeCursor queue.RingBufferCursor // concurrent write
	readCursor  uint64                 // synchronize read
	capacity    uint64
}

func NewXSequencer(capacity uint64) Sequencer {
	return &xSequencer{
		capacity:    capacity,
		readCursor:  1,
		writeCursor: queue.NewXRingBufferCursor(),
	}
}

func (x *xSequencer) NextReadCursor() uint64 {
	return atomic.AddUint64(&x.readCursor, 1)
}

func (x *xSequencer) LoadReadCursor() uint64 {
	return atomic.LoadUint64(&x.readCursor)
}

func (x *xSequencer) Capacity() uint64 {
	return atomic.LoadUint64(&x.capacity)
}

func (x *xSequencer) GetWriteCursor() queue.RingBufferCursor {
	return x.writeCursor
}
