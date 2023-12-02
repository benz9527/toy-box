package list

import (
	"sync"
	"sync/atomic"
)

var (
	_ BasicLinkedList[struct{}] = (*SinglyLinkedList[struct{}])(nil)           // Type check assertion
	_ BasicLinkedList[struct{}] = (*ConcurrentSinglyLinkedList[struct{}])(nil) // Type check assertion
)

type SinglyLinkedList[T comparable] struct {
	// sentinel list element.
	root nodeElement[T]
	len  atomic.Int64
}

func newSinglyLinkedList[T comparable](concurrent bool) BasicLinkedList[T] {
	l := new(SinglyLinkedList[T]).init()
	if concurrent {
		l.root.lock = &sync.RWMutex{}
	}
	return l
}

func NewSinglyLinkedList[T comparable]() BasicLinkedList[T] {
	return newSinglyLinkedList[T](false)
}

func (l *SinglyLinkedList[T]) getRoot() NodeElement[T] {
	return &l.root
}

func (l *SinglyLinkedList[T]) getRootHeader() NodeElement[T] {
	return l.root.next
}

func (l *SinglyLinkedList[T]) setRootHeader(targetE NodeElement[T]) {
	l.root.next = targetE
}

func (l *SinglyLinkedList[T]) getRootTail() NodeElement[T] {
	return l.root.prev
}

func (l *SinglyLinkedList[T]) setRootTail(targetE NodeElement[T]) {
	l.root.prev = targetE
}

func (l *SinglyLinkedList[T]) init() *SinglyLinkedList[T] {
	l.setRootHeader(&l.root)
	l.setRootTail(&l.root)
	l.root.list = l
	l.len.Store(0)
	return l
}

func (l *SinglyLinkedList[T]) Len() int64 {
	return l.len.Load()
}

func (l *SinglyLinkedList[T]) isListElement(targetE NodeElement[T]) bool {
	if targetE == nil {
		return false
	}
	target, ok := targetE.(*nodeElement[T])
	res := ok && target.hasLock() == l.root.hasLock()
	if res && target.list != nil {
		res = target.list == l
	}
	return res
}

func (l *SinglyLinkedList[T]) append(e NodeElement[T]) NodeElement[T] {
	if !l.isListElement(e) {
		return nil
	}
	if e.(*nodeElement[T]).list == nil {
		e.(*nodeElement[T]).list = l
	}
	if l.len.Load() <= 0 {
		l.setRootHeader(e)
	} else {
		l.getRootTail().(*nodeElement[T]).next = e
	}
	l.setRootTail(e)
	l.len.Add(1)
	return e
}

func (l *SinglyLinkedList[T]) Append(elements ...NodeElement[T]) []NodeElement[T] {
	for i := 0; i < len(elements); i++ {
		if elements[i] == nil {
			continue
		}
		e, ok := elements[i].(*nodeElement[T])
		if !ok {
			continue
		}
		elements[i] = l.append(e)
	}
	return elements
}

func (l *SinglyLinkedList[T]) AppendValue(values ...T) []NodeElement[T] {
	newElements := make([]NodeElement[T], 0, len(values))
	for _, v := range values {
		newElements = append(newElements, newNodeElement(v, l))
	}
	return l.Append(newElements...)
}

func (l *SinglyLinkedList[T]) InsertAfter(v T, dstE NodeElement[T]) NodeElement[T] {
	if !l.isListElement(dstE) {
		return nil
	}

	newE := newNodeElement[T](v, l)
	if dstE == l.getRoot() || dstE == l.getRootTail() {
		l.setRootTail(newE)
	}

	newE.next, dstE.(*nodeElement[T]).next = dstE.(*nodeElement[T]).next, newE
	l.len.Add(1)
	return newE
}

func (l *SinglyLinkedList[T]) InsertBefore(v T, dstE NodeElement[T]) NodeElement[T] {
	if !l.isListElement(dstE) {
		return nil
	}

	newE := newNodeElement[T](v, l)
	var iterator NodeElement[T] = nil
	if dstE == l.getRootHeader() {
		l.setRootHeader(newE)
	} else {
		iterator = l.getRoot()
		for iterator.HasNext() && iterator.GetNext() != dstE {
			iterator = iterator.GetNext()
		}
		// not found
		if iterator == nil {
			return nil
		}
		iterator.(*nodeElement[T]).next = newE
	}
	newE.next = dstE
	l.len.Add(1)
	return newE
}

func (l *SinglyLinkedList[T]) Remove(targetE NodeElement[T]) NodeElement[T] {
	if !l.isListElement(targetE) || targetE.(*nodeElement[T]).list == nil || l.len.Load() == 0 {
		return nil
	}

	defer func() {
		if l.len.Load() == 0 {
			l.setRootHeader(l.getRoot())
			l.setRootTail(l.getRoot())
		}
	}()

	var iterator = l.getRoot()
	// find previous element of targetE
	for iterator.HasNext() && iterator.GetNext() != targetE {
		iterator = iterator.GetNext()
	}
	// not found
	if iterator == nil {
		return nil
	}
	iterator.(*nodeElement[T]).next = targetE.GetNext()
	if targetE == l.getRootTail() {
		l.setRootTail(iterator)
	}
	l.len.Add(-1)
	return targetE
}

func (l *SinglyLinkedList[T]) ForEach(fn func(e NodeElement[T])) {
	if fn == nil || l.len.Load() == 0 || l.getRootHeader() == l.getRoot() {
		return
	}
	var iterator = l.getRoot()
	for iterator.HasNext() {
		fn(iterator.GetNext())
		iterator = iterator.GetNext()
	}
}

func (l *SinglyLinkedList[T]) FindFirst(targetV T, compareFn ...func(e NodeElement[T]) bool) (NodeElement[T], bool) {
	if l.len.Load() == 0 || l.getRootHeader() == l.getRoot() {
		return nil, false
	}

	if len(compareFn) <= 0 {
		compareFn = []func(e NodeElement[T]) bool{
			func(e NodeElement[T]) bool {
				return e.GetValue() == targetV
			},
		}
	}

	var iterator = l.getRoot()
	for iterator.HasNext() {
		if compareFn[0](iterator.GetNext()) {
			return iterator.GetNext(), true
		}
		iterator = iterator.GetNext()
	}
	return nil, false
}

type ConcurrentSinglyLinkedList[T comparable] struct {
	lock sync.RWMutex
	list BasicLinkedList[T]
}

func NewConcurrentSinglyLinkedList[T comparable]() BasicLinkedList[T] {
	slist := &ConcurrentSinglyLinkedList[T]{
		lock: sync.RWMutex{},
		list: newSinglyLinkedList[T](true),
	}
	return slist
}

func (l *ConcurrentSinglyLinkedList[T]) Len() int64 {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.list.Len()
}

func (l *ConcurrentSinglyLinkedList[T]) Append(elements ...NodeElement[T]) []NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	elements = l.list.Append(elements...)
	return elements
}

func (l *ConcurrentSinglyLinkedList[T]) AppendValue(values ...T) []NodeElement[T] {
	elements := make([]NodeElement[T], 0, len(values))
	for _, v := range values {
		elements = append(elements, newConcurrentNodeElement[T](v, l.list))
	}
	elements = l.Append(elements...)
	return elements
}

func (l *ConcurrentSinglyLinkedList[T]) InsertAfter(v T, dstE NodeElement[T]) NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	dstE = l.list.InsertAfter(v, dstE)
	return dstE
}

func (l *ConcurrentSinglyLinkedList[T]) InsertBefore(v T, dstE NodeElement[T]) NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	dstE = l.list.InsertBefore(v, dstE)
	return dstE
}

func (l *ConcurrentSinglyLinkedList[T]) Remove(targetE NodeElement[T]) NodeElement[T] {
	l.lock.Lock()
	defer l.lock.Unlock()
	targetE = l.list.Remove(targetE)
	return targetE
}

func (l *ConcurrentSinglyLinkedList[T]) ForEach(fn func(e NodeElement[T])) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.list.ForEach(fn)
}

func (l *ConcurrentSinglyLinkedList[T]) FindFirst(v T, compareFn ...func(e NodeElement[T]) bool) (NodeElement[T], bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	e, ok := l.list.FindFirst(v, compareFn...)
	return e, ok
}
