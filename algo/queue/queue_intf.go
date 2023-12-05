package queue

import "context"

// Here I defined the interfaces for priority queue and delay queue.
// In my plan, I will implement the priority queue and delay queue with array and heap.

type LessThan[E comparable, T PQItem[E]] func(i, j T) bool

type PQItem[E comparable] interface {
	GetPriority() int64
	SetPriority(priority int64)
	GetLessThan() LessThan[E, PQItem[E]]
	SetLessThan(comparator LessThan[E, PQItem[E]])
	GetIndex() int64
	SetIndex(index int64)
	GetValue() E
}

type PriorityQueue[E comparable, T PQItem[E]] interface {
	Len() int64
	Push(item T)
	Pop() (T, error)
	Peek() (T, error)
	PopAndShift(boundary int64) (T, int64)
}

type DQItem[E comparable] interface {
	PQItem[E]
	GetExpiration() int64
}

type DelayQueue[E comparable, T DQItem[E]] interface {
	Offer(item T, expiration int64) error
	Poll(ctx context.Context, nowFn func() int64) (T, error)
}
