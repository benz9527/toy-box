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

type PriorityItem[T any] struct {
	Value    T
	Priority int
	index    int
}

type _queue[T any] []*PriorityItem[T]

func (q *_queue[T]) Len() int {
	return len(*q)
}

func (q *_queue[T]) Less(i, j int) bool {
	return false
}

func (q *_queue[T]) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
	(*q)[i].index = i
	(*q)[j].index = j
}

func (q *_queue[T]) Pop() interface{} {
	prev := *q
	n := len(prev)
	item := prev[n-1]
	item.index = -1
	prev[n-1] = nil
	*q = prev[:n-1]
	return item
}

func (q *_queue[T]) Push(i interface{}) {
	item := i.(*PriorityItem[T])
	prev := *q
	n := len(prev)
	item.index = n
	*q = append(*q, item)
}

type PriorityQueue[T any] struct {
	data _queue[T]
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		data: make(_queue[T], 0, 32),
	}
	heap.Init(&pq.data)
	return pq
}

func (pq *PriorityQueue[T]) Pop() (T, error) {
	if len(pq.data) == 0 {
		nilT := new(T)
		return *nilT, errors.New("empty")
	}
	item := heap.Pop(&pq.data)
	return item.(T), nil
}

func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(&pq.data, item)
}
