package concurrency

import (
	"fmt"
	"sync"
)

type GenericChannel[T comparable] struct {
	ch       chan T
	isClosed bool
	once     *sync.Once
	rwLock   *sync.RWMutex
}

func NewGenericChannel[T comparable](chSize ...int) *GenericChannel[T] {
	isNoCacheCh := true
	size := 1
	if len(chSize) > 0 {
		if chSize[0] > 0 {
			size = chSize[0]
			isNoCacheCh = false
		}
	}
	if isNoCacheCh {
		ch := make(chan T)
		return &GenericChannel[T]{ch: ch, once: &sync.Once{}, rwLock: &sync.RWMutex{}}
	}
	ch := make(chan T, size)
	return &GenericChannel[T]{ch: ch, once: &sync.Once{}, rwLock: &sync.RWMutex{}}
}

func (c *GenericChannel[T]) IsClosed() bool {
	return c.isClosed
}

func (c *GenericChannel[T]) Close() {
	c.once.Do(func() {
		c.rwLock.Lock()
		c.isClosed = true
		close(c.ch)
		c.rwLock.Unlock()
	})
}

func (c *GenericChannel[T]) Wait() <-chan T {
	return c.ch
}

func (c *GenericChannel[T]) Send(v T) error {
	c.rwLock.RLock()
	if c.IsClosed() {
		c.rwLock.RUnlock()
		return fmt.Errorf("channel has been closed")
	}
	c.rwLock.RUnlock()

	c.ch <- v
	return nil
}
