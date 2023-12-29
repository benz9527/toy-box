package pubsub

import (
	"github.com/benz9527/toy-box/algo/queue"
)

type xSequencer struct {
	writeCursor queue.RingBufferCursor // concurrent write
	readCursor  queue.RingBufferCursor // concurrent read
	capacity    uint64
}

func NewXSequencer(capacity uint64) Sequencer {
	return &xSequencer{
		capacity:    capacity,
		readCursor:  queue.NewXRingBufferCursor(),
		writeCursor: queue.NewXRingBufferCursor(),
	}
}

func (x *xSequencer) GetReadCursor() queue.RingBufferCursor {
	return x.readCursor
}

func (x *xSequencer) Capacity() uint64 {
	return x.capacity
}

func (x *xSequencer) GetWriteCursor() queue.RingBufferCursor {
	return x.writeCursor
}

type xSequenceReconciler struct {
}
