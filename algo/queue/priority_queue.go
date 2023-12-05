package queue

import (
	"container/heap"
	"errors"
)

// 使用最小堆来构建优先级队列
// 每次最小堆的根节点都会与树中最后一个叶子节点交换位置
// 之后剪去该叶子节点就获得了当前优先级最高的信息
// 剪枝的操作完成后，原本被提升为根节点的节点会被重新下沉
// 完成优先级的重新分布
// 新加入节点也会重新排列分布

type pqItem[T comparable] struct {
	value      T
	priority   int64
	index      int64
	comparator LessThan[T, PQItem[T]]
}

func (item *pqItem[T]) GetPriority() int64 {
	return item.priority
}

func (item *pqItem[T]) SetPriority(priority int64) {
	item.priority = priority
}

func (item *pqItem[T]) GetIndex() int64 {
	return item.index
}

func (item *pqItem[T]) SetIndex(index int64) {
	item.index = index
}

func (item *pqItem[T]) GetValue() T {
	return item.value
}

func (item *pqItem[T]) GetLessThan() LessThan[T, PQItem[T]] {
	return item.comparator
}

func (item *pqItem[T]) SetLessThan(comparator LessThan[T, PQItem[T]]) {
	if comparator == nil {
		// default comparator, default min priority value has high priority
		comparator = func(i, j PQItem[T]) bool {
			return i.GetPriority() < j.GetPriority()
		}
	}
	item.comparator = comparator
}

func NewPQItem[T comparable](value T, priority int64) PQItem[T] {
	return &pqItem[T]{
		value:    value,
		priority: priority,
		index:    0,
		comparator: func(i, j PQItem[T]) bool {
			// default comparator, default min priority value has high priority
			return i.GetPriority() < j.GetPriority()
		},
	}
}

type arrayPQ[V comparable, T PQItem[V]] []T

func (q *arrayPQ[V, T]) Len() int {
	return len(*q)
}

func (q *arrayPQ[V, T]) Less(i, j int) bool {
	// Nil panic at your own risk
	return (*q)[i].GetLessThan()((*q)[i], (*q)[j])
}

func (q *arrayPQ[V, T]) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
	(*q)[i].SetIndex(int64(i))
	(*q)[j].SetIndex(int64(j))
}

func (q *arrayPQ[V, T]) Pop() interface{} {
	prev := *q
	n := len(prev)
	item := prev[n-1]
	item.SetIndex(-1)
	prev[n-1] = *new(T) // nil object
	*q = prev[:n-1]
	return item
}

func (q *arrayPQ[V, T]) Push(i interface{}) {
	item, ok := i.(PQItem[V])
	if !ok {
		return
	}

	prev := *q
	n := len(prev)
	item.SetIndex(int64(n))
	*q = append(*q, item.(T))
}

type ArrayPriorityQueue[V comparable, T PQItem[V]] struct {
	queue           arrayPQ[V, T]
	localComparator LessThan[V, PQItem[V]]
}

// NewArrayPriorityQueue create a priority queue, you can customize the less than comparator.
// Default less than comparator implemented by min priority value has high priority.
func NewArrayPriorityQueue[V comparable, T PQItem[V]](
	capacity int,
	comparator ...LessThan[V, PQItem[V]],
) PriorityQueue[V, T] {
	if len(comparator) <= 0 {
		comparator = []LessThan[V, PQItem[V]]{
			func(i, j PQItem[V]) bool {
				return i.GetPriority() < j.GetPriority()
			},
		}
	}
	if capacity <= 0 {
		capacity = 32
	}
	pq := &ArrayPriorityQueue[V, T]{
		queue:           make(arrayPQ[V, T], 0, capacity),
		localComparator: comparator[0],
	}
	heap.Init(&pq.queue)
	return pq
}

func (pq *ArrayPriorityQueue[V, T]) Len() int64 {
	return int64(len(pq.queue))
}

func (pq *ArrayPriorityQueue[V, T]) Pop() (T, error) {
	if len(pq.queue) == 0 {
		nilT := new(T)
		return *nilT, errors.New("empty")
	}
	item := heap.Pop(&pq.queue)
	return item.(T), nil
}

func (pq *ArrayPriorityQueue[V, T]) Push(item T) {
	item.SetLessThan(pq.localComparator)
	heap.Push(&pq.queue, item)
}

func (pq *ArrayPriorityQueue[V, T]) Peek() (T, error) {
	if len(pq.queue) == 0 {
		nilT := new(T)
		return *nilT, errors.New("empty")
	}
	return pq.queue[0], nil
}

func (pq *ArrayPriorityQueue[V, T]) PopAndShift(boundary int64) (T, int64) {
	nilT := new(T)
	if (*pq).queue.Len() == 0 {
		return *nilT, 0
	}
	item := (*pq).queue[0]
	if item.GetPriority() > boundary {
		return *nilT, item.GetPriority() - boundary
	}
	heap.Remove(&(*pq).queue, 0)
	return item, 0
}
