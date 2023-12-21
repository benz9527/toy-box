package pubsub

import (
	"errors"
	"github.com/benz9527/toy-box/algo/queue"
	"sync/atomic"
)

type xSequenceBarrier struct {
	alerted          *atomic.Bool
	strategy         BlockStrategy
	producerSeq      ProducerSequencer
	dependentCursors []queue.RingBufferCursor
}

func (barrier *xSequenceBarrier) CheckAlert() error {
	if barrier.alerted.Load() {
		return errors.New("alerted")
	}
	return nil
}

func (barrier *xSequenceBarrier) ClearAlert() {
	barrier.alerted.Store(false)
}

func (barrier *xSequenceBarrier) Alert() {
	barrier.alerted.Store(true)
	barrier.strategy.Done()
}

func (barrier *xSequenceBarrier) GetAvailableConsumeCursor(currentConsumeCursor uint64) uint64 {
	if err := barrier.CheckAlert(); err != nil {
		// Producer should be stopped for a while
	}

}
