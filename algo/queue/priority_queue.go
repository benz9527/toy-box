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

type pqItem[E comparable] struct {
	value      E
	priority   int64
	index      int64
	comparator LessThan[E]
}

func (item *pqItem[E]) GetPriority() int64 {
	return item.priority
}

func (item *pqItem[E]) SetPriority(priority int64) {
	item.priority = priority
}

func (item *pqItem[E]) GetIndex() int64 {
	return item.index
}

func (item *pqItem[E]) SetIndex(index int64) {
	item.index = index
}

func (item *pqItem[E]) GetValue() E {
	return item.value
}

func (item *pqItem[E]) GetLessThan() LessThan[E] {
	return item.comparator
}

func (item *pqItem[E]) SetLessThan(comparator LessThan[E]) {
	if comparator == nil {
		// default comparator, default min priority value has high priority
		comparator = func(i, j PQItem[E]) bool {
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

type arrayPQ[E comparable] []PQItem[E]

func (q *arrayPQ[E]) Len() int {
	return len(*q)
}

func (q *arrayPQ[E]) Less(i, j int) bool {
	// Nil panic at your own risk
	return (*q)[i].GetLessThan()((*q)[i], (*q)[j])
}

func (q *arrayPQ[E]) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
	(*q)[i].SetIndex(int64(i))
	(*q)[j].SetIndex(int64(j))
}

func (q *arrayPQ[E]) Pop() interface{} {
	prev := *q
	n := len(prev)
	item := prev[n-1]
	item.SetIndex(-1)
	prev[n-1] = *new(PQItem[E]) // nil object
	*q = prev[:n-1]
	return item
}

func (q *arrayPQ[E]) Push(i interface{}) {
	item, ok := i.(PQItem[E])
	if !ok {
		return
	}

	prev := *q
	n := len(prev)
	item.SetIndex(int64(n))
	*q = append(*q, item.(PQItem[E]))
}

type ArrayPriorityQueue[E comparable] struct {
	queue           arrayPQ[E]
	localComparator LessThan[E]
}

// NewArrayPriorityQueue create a priority queue, you can customize the less than comparator.
// Default less than comparator implemented by min priority value has high priority.
func NewArrayPriorityQueue[E comparable](
	capacity int,
	comparator ...LessThan[E],
) PriorityQueue[E] {
	if len(comparator) <= 0 {
		comparator = []LessThan[E]{
			func(i, j PQItem[E]) bool {
				return i.GetPriority() < j.GetPriority()
			},
		}
	}
	if capacity <= 0 {
		capacity = 32
	}
	pq := &ArrayPriorityQueue[E]{
		queue:           make(arrayPQ[E], 0, capacity),
		localComparator: comparator[0],
	}
	heap.Init(&pq.queue)
	return pq
}

func (pq *ArrayPriorityQueue[E]) Len() int64 {
	return int64(len(pq.queue))
}

func (pq *ArrayPriorityQueue[E]) Pop() (PQItem[E], error) {
	if len(pq.queue) == 0 {
		nilT := new(PQItem[E])
		return *nilT, errors.New("empty")
	}
	item := heap.Pop(&pq.queue)
	return item.(PQItem[E]), nil
}

func (pq *ArrayPriorityQueue[E]) Push(item PQItem[E]) {
	item.SetLessThan(pq.localComparator)
	heap.Push(&pq.queue, item)
}

func (pq *ArrayPriorityQueue[E]) Peek() (PQItem[E], error) {
	if len(pq.queue) == 0 {
		nilT := new(PQItem[E])
		return *nilT, errors.New("empty")
	}
	return pq.queue[0], nil
}

func (pq *ArrayPriorityQueue[E]) PopIfMatched(boundary int64) (PQItem[E], int64) {
	nilT := new(PQItem[E])
	if (*pq).queue.Len() == 0 {
		return *nilT, 0
	}
	item := (*pq).queue[0]
	if item.GetPriority() > boundary {
		// not matched
		return *nilT, item.GetPriority() - boundary
	}
	heap.Remove(&(*pq).queue, 0)
	return item, 0
}
