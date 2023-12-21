package pubsub

import (
	"github.com/benz9527/toy-box/algo/queue"
	"math"
	"runtime"
	"sync/atomic"
	"unsafe"
)

type xBasicSequencer struct {
	writeCursor queue.RingBufferCursor // concurrent write
	readCursor  queue.RingBufferCursor // concurrent read
	capacity    uint64
}

func NewXBasicSequencer(capacity uint64) Sequencer {
	return &xBasicSequencer{
		capacity:    capacity,
		readCursor:  queue.NewXRingBufferCursor(),
		writeCursor: queue.NewXRingBufferCursor(),
	}
}

func (x *xBasicSequencer) GetReadCursor() queue.RingBufferCursor {
	return x.readCursor
}

func (x *xBasicSequencer) Capacity() uint64 {
	return x.capacity
}

func (x *xBasicSequencer) GetWriteCursor() queue.RingBufferCursor {
	return x.writeCursor
}

// Same as Java Disruptor.SingleProducerSequencer
// Single goroutine producer
type xSingleProducerSequencer struct {
	// As Java volatile by AtomicReferenceFieldUpdater<xSingleProducerSequencer, queue.RingBufferCursor[]>
	gatingConsumerCursors unsafe.Pointer // []queue.RingBufferCursor
	currentProducerCursor queue.RingBufferCursor
	strategy              BlockStrategy
	capacity              uint64 // ring buffer capacity
	// FIXME notice the false sharing data race
	nextProducerCursorValue   uint64
	cachedConsumerCursorValue uint64
}

func newXSingleProducerSequencer(
	capacity uint64,
	strategy BlockStrategy,
) *xSingleProducerSequencer {
	return &xSingleProducerSequencer{}
}

func (seq *xSingleProducerSequencer) getGatingConsumerCursors() []queue.RingBufferCursor {
	return *(*[]queue.RingBufferCursor)(atomic.LoadPointer(&seq.gatingConsumerCursors))
}

func (seq *xSingleProducerSequencer) NextN(n uint64) uint64 {
	// Batch request producer next n slots
	nextProducerNCursorVal := seq.nextProducerCursorValue + n
	// Producer cursor value has been greater than consumer cursor value
	// a ring buffer capacity
	wrapPoint := nextProducerNCursorVal - seq.capacity
	cachedGatingConsumerCursorVal := seq.cachedConsumerCursorValue
	if wrapPoint > cachedGatingConsumerCursorVal {
		// Check consumer cursors progress
		minCursorVal := getMinCursorValueBy(seq.getGatingConsumerCursors(), nextProducerNCursorVal)
		for ; wrapPoint > minCursorVal; minCursorVal = getMinCursorValueBy(seq.getGatingConsumerCursors(), nextProducerNCursorVal) {
			runtime.Gosched()
		}
		seq.cachedConsumerCursorValue = minCursorVal
	}
	seq.nextProducerCursorValue = nextProducerNCursorVal
	return nextProducerNCursorVal
}

func (seq *xSingleProducerSequencer) Next() uint64 {
	return seq.NextN(1)
}

func (seq *xSingleProducerSequencer) Update(publishedCursor uint64) {
	seq.currentProducerCursor.CASStore(seq.currentProducerCursor.AtomicLoad(), publishedCursor)
	// FIXME Notify blocking consumer by blocked strategy
	seq.strategy.Done()
}

func (seq *xSingleProducerSequencer) AddGatingConsumerCursors(cursors ...queue.RingBufferCursor) {
	// concurrent add gating consumer cursors
	for {
		// Spin and hotspot code
		updatedCursors := make([]queue.RingBufferCursor, 0, len(seq.getGatingConsumerCursors())+len(cursors))
		copy(updatedCursors, seq.getGatingConsumerCursors())
		currentProducerCursorVal := seq.currentProducerCursor.AtomicLoad()
		for _, cursor := range cursors {
			cursor.CASStore(cursor.AtomicLoad(), currentProducerCursorVal)
			updatedCursors = append(updatedCursors, cursor)
		}
		oldRef := seq.gatingConsumerCursors
		if atomic.CompareAndSwapPointer(&seq.gatingConsumerCursors, oldRef, unsafe.Pointer(&updatedCursors)) {
			break
		}
	}

	// Update again
	currentProducerCursorVal := seq.currentProducerCursor.AtomicLoad()
	for _, cursor := range cursors {
		cursor.CASStore(cursor.AtomicLoad(), currentProducerCursorVal)
	}
}

func (seq *xSingleProducerSequencer) RemoveGatingConsumerCursor(cursor queue.RingBufferCursor) bool {
	// concurrent remove gating consumer cursors
	counting := func(oldCursors []queue.RingBufferCursor, rmTarget queue.RingBufferCursor) int {
		numToRemove := 0
		for _, cursor := range oldCursors {
			if cursor == rmTarget {
				numToRemove++
			}
		}
		return numToRemove
	}
	numToRemove := 0
	for {
		// Spin and hotspot code
		oldCursors := seq.getGatingConsumerCursors()
		numToRemove = counting(oldCursors, cursor)
		if numToRemove == 0 {
			break
		}
		updatedCursors := make([]queue.RingBufferCursor, 0, len(seq.getGatingConsumerCursors())-numToRemove)
		for _, oldCursor := range oldCursors {
			if oldCursor != cursor {
				oldCursor.CASStore(cursor.AtomicLoad(), math.MaxUint64)
			}
		}
		oldRef := seq.gatingConsumerCursors
		if atomic.CompareAndSwapPointer(&seq.gatingConsumerCursors, oldRef, unsafe.Pointer(&updatedCursors)) {
			break
		}
	}
	return numToRemove != 0
}

type xMultiProducerSequencer struct {
}

// Sequence it is the cursor value
func getMinCursorValueBy(cursors []queue.RingBufferCursor, minCursorVal uint64) uint64 {
	for _, cursor := range cursors {
		cursorVal := cursor.AtomicLoad()
		// get min cursor value
		if cursorVal < minCursorVal {
			minCursorVal = cursorVal
		}
	}
	return minCursorVal
}

func getMinCursorVale(cursors []queue.RingBufferCursor) uint64 {
	return getMinCursorValueBy(cursors, math.MaxUint64)
}
