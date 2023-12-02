package list

import (
	"sync"
	"sync/atomic"
)

type NodeElement[T comparable] interface {
	GetNext() NodeElement[T]
	GetPrev() NodeElement[T]
	GetValue() T
}

type List[T comparable] interface {
	Len() int64
	Append(e NodeElement[T]) NodeElement[T]
	AppendValue(v T) NodeElement[T]
	InsertAfter(e NodeElement[T], v T) NodeElement[T]
	InsertBefore(e NodeElement[T], v T) NodeElement[T]
	Remove(e NodeElement[T]) NodeElement[T]
	ForEach(fn func(e NodeElement[T]))
	Find(v T, compareFn ...func(e NodeElement[T]) bool) (NodeElement[T], bool)
}

var (
	_ NodeElement[struct{}] = (*SinglyNodeElement[struct{}])(nil)          // Type check assertion
	_ List[struct{}]        = (*SinglyLinkedList[struct{}])(nil)           // Type check assertion
	_ List[struct{}]        = (*ConcurrentSinglyLinkedList[struct{}])(nil) // Type check assertion
)

type SinglyNodeElement[T comparable] struct {
	next  NodeElement[T]
	list  List[T]
	value T
}

func NewSinglyNodeElement[T comparable](v T) NodeElement[T] {
	return &SinglyNodeElement[T]{
		value: v,
	}
}

func newSinglyNodeElement[T comparable](v T, list List[T]) *SinglyNodeElement[T] {
	return &SinglyNodeElement[T]{
		value: v,
		list:  list,
	}
}

func (e *SinglyNodeElement[T]) GetNext() NodeElement[T] {
	n := e.next
	if n == nil {
		return nil
	}
	nse, ok := n.(*SinglyNodeElement[T])
	if !ok {
		return nil
	}
	if nse.list != e.list {
		return nil
	}
	return n
}

func (e *SinglyNodeElement[T]) GetPrev() NodeElement[T] {
	return nil
}

func (e *SinglyNodeElement[T]) GetValue() T {
	return e.value
}

type SinglyLinkedList[T comparable] struct {
	// sentinel list element, only &root and root.next are used.
	// This is the reason why it is not a pointer reference.
	root SinglyNodeElement[T]
	tail *SinglyNodeElement[T]
	len  atomic.Int64
}

func NewSinglyLinkedList[T comparable]() List[T] {
	return new(SinglyLinkedList[T]).init()
}

func (l *SinglyLinkedList[T]) init() *SinglyLinkedList[T] {
	l.root.next = &l.root
	l.root.list = l
	l.tail = &l.root
	l.tail.list = l
	l.len.Store(0)
	return l
}

func (l *SinglyLinkedList[T]) Len() int64 {
	return l.len.Load()
}

func (l *SinglyLinkedList[T]) Append(e NodeElement[T]) NodeElement[T] {
	sne, ok := e.(*SinglyNodeElement[T])
	if !ok {
		return l.tail
	}
	sne.list = l
	if l.root.next == &l.root {
		l.root.next = sne
	}
	if l.tail != &l.root {
		pe := l.tail
		pe.next = sne
	}
	l.tail = sne
	l.len.Add(1)
	return sne
}

func (l *SinglyLinkedList[T]) AppendValue(v T) NodeElement[T] {
	return l.Append(NewSinglyNodeElement[T](v))
}

func (l *SinglyLinkedList[T]) InsertAfter(e NodeElement[T], v T) NodeElement[T] {
	if e == nil {
		return l.tail
	}

	psne, ok := e.(*SinglyNodeElement[T])
	if !ok {
		return l.tail
	}
	if psne.list != l {
		return l.tail
	}

	sne := newSinglyNodeElement[T](v, l)
	if psne == l.tail {
		l.tail = sne
	}
	n := psne.next
	psne.next = sne
	sne.next = n
	l.len.Add(1)
	return sne
}

func (l *SinglyLinkedList[T]) InsertBefore(e NodeElement[T], v T) NodeElement[T] {
	if e == nil {
		return l.tail
	}

	nsne, ok := e.(*SinglyNodeElement[T])
	if !ok {
		return l.tail
	}
	if nsne.list != l {
		return l.tail
	}
	sne := newSinglyNodeElement[T](v, l)
	var pe NodeElement[T] = nil
	if e == l.root.next {
		l.root.next = sne
	} else {
		pe = &l.root
		for pe.GetNext() != e {
			pe = pe.GetNext()
		}
		pe.(*SinglyNodeElement[T]).next = sne
	}
	sne.next = nsne
	l.len.Add(1)
	return sne
}

func (l *SinglyLinkedList[T]) Remove(e NodeElement[T]) NodeElement[T] {
	if e == nil {
		return nil
	}
	var pe = l.root.next
	if pe == nil || pe == &l.root {
		return nil
	} else if pe == e {
		l.root.next = pe.GetNext()
	} else {
		for pe.GetNext() != e {
			pe = pe.GetNext()
		}
		pe.(*SinglyNodeElement[T]).next = e.GetNext()
	}
	if e == l.tail {
		l.tail = pe.(*SinglyNodeElement[T])
	}
	l.len.Add(-1)
	return e
}

func (l *SinglyLinkedList[T]) ForEach(fn func(e NodeElement[T])) {
	var iterator NodeElement[T] = &l.root
	if iterator == nil || fn == nil || l.len.Load() == 0 || l.root.next == &l.root {
		return
	}
	for iterator.GetNext() != nil {
		fn(iterator.GetNext())
		iterator = iterator.GetNext()
	}
}

func (l *SinglyLinkedList[T]) Find(v T, compareFn ...func(e NodeElement[T]) bool) (NodeElement[T], bool) {
	var iterator NodeElement[T] = &l.root
	if iterator == nil || l.len.Load() == 0 || l.root.next == &l.root {
		return nil, false
	}

	if len(compareFn) <= 0 {
		compareFn = []func(e NodeElement[T]) bool{
			func(e NodeElement[T]) bool {
				return e.GetValue() == v
			},
		}
	}

	for iterator.GetNext() != nil {
		if compareFn[0](iterator.GetNext()) {
			return iterator.GetNext(), true
		}
		iterator = iterator.GetNext()
	}
	return nil, false
}

type ConcurrentSinglyLinkedList[T comparable] struct {
	lock sync.RWMutex
	list List[T]
}

func NewConcurrentSinglyLinkedList[T comparable]() List[T] {
	slist := &ConcurrentSinglyLinkedList[T]{
		lock: sync.RWMutex{},
		list: NewSinglyLinkedList[T](),
	}
	return slist
}

func (l *ConcurrentSinglyLinkedList[T]) Len() int64 {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.list.Len()
}

func (l *ConcurrentSinglyLinkedList[T]) Append(e NodeElement[T]) NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	e = l.list.Append(e)
	return e
}

func (l *ConcurrentSinglyLinkedList[T]) AppendValue(v T) NodeElement[T] {
	return l.Append(NewSinglyNodeElement[T](v))
}

func (l *ConcurrentSinglyLinkedList[T]) InsertAfter(e NodeElement[T], v T) NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	e = l.list.InsertAfter(e, v)
	return e
}

func (l *ConcurrentSinglyLinkedList[T]) InsertBefore(e NodeElement[T], v T) NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	e = l.list.InsertBefore(e, v)
	return e
}

func (l *ConcurrentSinglyLinkedList[T]) Remove(e NodeElement[T]) NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	e = l.list.Remove(e)
	return e
}

func (l *ConcurrentSinglyLinkedList[T]) ForEach(fn func(e NodeElement[T])) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.list.ForEach(fn)
}

func (l *ConcurrentSinglyLinkedList[T]) Find(v T, compareFn ...func(e NodeElement[T]) bool) (NodeElement[T], bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	e, ok := l.list.Find(v, compareFn...)
	return e, ok
}
