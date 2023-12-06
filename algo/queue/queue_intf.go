package queue

import "context"

// Here I defined the interfaces for priority queue and delay queue.
// In my plan, I will implement the priority queue and delay queue with array and heap.

type LessThan[E comparable] func(i, j PQItem[E]) bool

type PQItem[E comparable] interface {
	GetPriority() int64
	SetPriority(priority int64)
	GetLessThan() LessThan[E]
	SetLessThan(comparator LessThan[E])
	GetIndex() int64
	SetIndex(index int64)
	GetValue() E
}

type PriorityQueue[E comparable] interface {
	Len() int64
	Push(item PQItem[E])
	Pop() PQItem[E]
	Peek() PQItem[E]
}

type DQItem[E comparable] interface {
	GetExpiration() int64
	GetPQItem() PQItem[E]
}

type DelayQueue[E comparable] interface {
	Offer(item E, expiration int64) error
	Poll(ctx context.Context, nowFn func() int64)
	Wait() <-chan E
}
