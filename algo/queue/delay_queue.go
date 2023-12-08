package queue

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errEmptyContext = errors.New("empty context")
)

type dqItem[E comparable] struct {
	PQItem[E]
}

func NewDQItem[E comparable](value E, expiration int64) DQItem[E] {
	return &dqItem[E]{
		PQItem: NewPQItem[E](value, expiration),
	}
}

func (d *dqItem[E]) GetPQItem() PQItem[E] {
	return d
}

func (d *dqItem[E]) GetExpiration() int64 {
	return d.GetPriority()
}

type sleepEnum = int32

const (
	wakeUp sleepEnum = iota
	fallAsleep
)

// size: 48
type arrayDQ[E comparable] struct {
	pq       PriorityQueue[E] // alignment size: 8; size: 16
	wakeUpC  chan struct{}    // alignment size: 8; size: 8
	lock     *sync.RWMutex    // alignment size: 8; size: 8
	mu       *sync.Mutex      // alignment size: 8; size: 8
	sleeping int32            // alignment size: 4; size: 4
}

func NewArrayDelayQueue[E comparable](
	capacity int,
	comparator ...LessThan[E],
) DelayQueue[E] {
	if capacity <= 0 {
		capacity = 32
	}
	if len(comparator) <= 0 {
		comparator = []LessThan[E]{
			func(i, j PQItem[E]) bool {
				return i.GetPriority() < j.GetPriority()
			},
		}
	}
	return &arrayDQ[E]{
		wakeUpC: make(chan struct{}, 1),
		pq:      NewArrayPriorityQueue[E](capacity, comparator[0]),
		mu:      &sync.Mutex{},
		lock:    &sync.RWMutex{},
	}
}

func (dq *arrayDQ[E]) popIfExpired(expiredBoundary int64) (item PQItem[E], deltaMs int64) {
	if (*dq).pq.Len() == 0 {
		return nil, 0
	}

	item = (*dq).pq.Peek()
	expiration := item.(DQItem[E]).GetExpiration()
	if expiration > expiredBoundary {
		// not matched
		return nil, expiration - expiredBoundary
	}
	item = (*dq).pq.Pop()
	return item, 0
}

func (dq *arrayDQ[E]) Offer(item E, expiration int64) error {
	e := NewDQItem[E](item, expiration)
	dq.lock.Lock()
	dq.pq.Push(e.GetPQItem())
	dq.lock.Unlock()

	if e.GetPQItem().GetIndex() == 0 {
		// Highest priority item, wake up the consumer
		if atomic.CompareAndSwapInt32(&dq.sleeping, fallAsleep, wakeUp) {
			dq.wakeUpC <- struct{}{}
		}
	}
	return nil
}

func (dq *arrayDQ[E]) poll(ctx context.Context, nowFn func() int64, sender chan<- E, closeChAfterFinish bool) {
	dq.mu.Lock() // Avoid concurrent execution of Poll()
	defer func() {
		// FIXME recover defer execution order
		if err := recover(); err != nil {
			slog.Error("delay queue panic recover", "error", err)
		}
		// before exit
		atomic.StoreInt32(&dq.sleeping, wakeUp)
		if closeChAfterFinish && sender != nil {
			close(sender)
		}
		dq.mu.Unlock()
	}()
	for {
		now := nowFn()
		dq.lock.RLock() // Concurrency control, avoid long time lock, block Offer()
		item, deltaMs := dq.popIfExpired(now)
		if item == nil {
			// No expired item in the queue
			// 1. without any item in the queue
			// 2. all items in the queue are not expired
			atomic.StoreInt32(&dq.sleeping, fallAsleep)
		}
		dq.lock.RUnlock()

		if item == nil {
			if deltaMs == 0 {
				// Queue is empty, waiting for new item
				select {
				case <-ctx.Done():
					return
				case <-dq.wakeUpC:
					// Waiting for an immediately executed item
					continue
				}
			} else if deltaMs > 0 {
				select {
				case <-ctx.Done():
					return
				case <-dq.wakeUpC:
					continue
				case <-time.After(time.Duration(deltaMs) * time.Millisecond):
					// Waiting for this item to be expired
					if atomic.SwapInt32(&dq.sleeping, wakeUp) == fallAsleep {
						dq.wakeUpC <- struct{}{}
					}
					continue
				}
			}
		}

		select {
		case <-ctx.Done():
			return
		case sender <- item.GetValue():
			// Waiting for the consumer to consume this item
			// If external channel is closed, here will be panic
		}
	}
}

func (dq *arrayDQ[E]) Poll(ctx context.Context, nowFn func() int64) (<-chan E, error) {
	if ctx == nil {
		return nil, errEmptyContext
	}
	resultC := make(chan E)
	// FIXME using goroutine pool
	go dq.poll(ctx, nowFn, resultC, true)
	return resultC, nil
}

func (dq *arrayDQ[E]) PollToChannel(ctx context.Context, nowFn func() int64, C chan<- E) error {
	if ctx == nil {
		return errEmptyContext
	}
	dq.poll(ctx, nowFn, C, false)
	return nil
}
