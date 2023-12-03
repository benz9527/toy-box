package list

import (
	"sync/atomic"
)

var (
	_ LinkedList[struct{}] = (*doublyLinkedList[struct{}])(nil) // Type check assertion
)

type nodeElementInListStatus uint8

const (
	notInList nodeElementInListStatus = iota
	theOnlyOne
	theFirstButNotTheLast
	theLastButNotTheFirst
	inMiddle
)

type doublyLinkedList[T comparable] struct {
	root nodeElement[T]
	len  atomic.Int64
}

func NewLinkedList[T comparable]() LinkedList[T] {
	return new(doublyLinkedList[T]).init()
}

func (l *doublyLinkedList[T]) getRoot() NodeElement[T] {
	return &l.root
}

func (l *doublyLinkedList[T]) getRootHeader() NodeElement[T] {
	return l.root.next
}

func (l *doublyLinkedList[T]) setRootHeader(targetE NodeElement[T]) {
	l.root.next = targetE
}

func (l *doublyLinkedList[T]) getRootTail() NodeElement[T] {
	return l.root.prev
}

func (l *doublyLinkedList[T]) setRootTail(targetE NodeElement[T]) {
	l.root.prev = targetE
}

func (l *doublyLinkedList[T]) init() *doublyLinkedList[T] {
	l.root.list = l
	l.setRootHeader(&l.root)
	l.setRootTail(&l.root)
	l.len.Store(0)
	return l
}

func (l *doublyLinkedList[T]) lazyInit() {
	if l.getRootHeader() == nil {
		l.init()
	}
}

func (l *doublyLinkedList[T]) Len() int64 {
	return l.len.Load()
}

func (l *doublyLinkedList[T]) checkElement(targetE NodeElement[T]) (*nodeElement[T], nodeElementInListStatus) {
	if targetE == nil {
		return nil, notInList
	}
	at, ok := targetE.(*nodeElement[T])
	if !ok || l.len.Load() == 0 || at.list != l || at.hasLock() != l.root.hasLock() {
		return nil, notInList
	}

	// mem address compare
	switch {
	case at.GetPrev() == nil && at.GetNext() == nil:
		// targetE is the first one and the last one
		if l.getRootHeader() != targetE || l.getRootTail() != targetE {
			return nil, notInList
		}
		return at, theOnlyOne
	case at.GetPrev() == nil && at.GetNext() != nil:
		// targetE is the first one but not the last one
		if at.GetNext() != targetE {
			return nil, notInList
		}
		// Ignore l.setRootTail (tail)
		return at, theFirstButNotTheLast
	case at.GetPrev() != nil && at.GetNext() == nil:
		// targetE is the last one but not the first one
		if at.GetPrev() != targetE {
			return nil, notInList
		}
		// Ignore l.setRootHeader (head)
		return at, theLastButNotTheFirst
	case at.GetPrev() != nil && at.GetNext() != nil:
		// targetE is neither the first one nor the last one
		if at.GetPrev() != targetE && at.GetNext() != targetE {
			return nil, notInList
		}
		return at, inMiddle
	}
	return nil, notInList
}

func (l *doublyLinkedList[T]) append(e *nodeElement[T]) NodeElement[T] {
	e.list = l
	e.next, e.prev = nil, nil

	if l.len.Load() == 0 {
		// empty list, new append element is the first one
		l.setRootHeader(e)
		l.setRootTail(e)
		return e
	}

	if l.root.GetPrev() == nil {
		return nil
	}
	lastOne := l.getRootTail().(*nodeElement[T])
	lastOne.next = e
	e.prev, e.next = lastOne, nil

	l.len.Add(1)
	return e
}

func (l *doublyLinkedList[T]) Append(elements ...NodeElement[T]) []NodeElement[T] {
	for i := 0; i < len(elements); i++ {
		if elements[i] == nil {
			continue
		}
		e, ok := elements[i].(*nodeElement[T])
		if !ok || e.list != l || l.root.hasLock() != e.hasLock() {
			continue
		}
		elements[i] = l.append(e)
	}
	return elements
}

func (l *doublyLinkedList[T]) AppendValue(values ...T) []NodeElement[T] {
	newElements := make([]NodeElement[T], 0, len(values))
	for _, v := range values {
		newElements = append(newElements, newNodeElement(v, l))
	}
	return l.Append(newElements...)
}

func (l *doublyLinkedList[T]) insertAfter(v T, at *nodeElement[T]) *nodeElement[T] {
	newE := newNodeElement(v, l)
	newE.prev = at
	newE.next = at.GetNext()
	at.next = newE
	if newE.GetNext() != nil {
		newE.GetNext().(*nodeElement[T]).prev = newE
	}

	l.len.Add(1)
	return newE
}

func (l *doublyLinkedList[T]) InsertAfter(v T, dstE NodeElement[T]) NodeElement[T] {
	at, status := l.checkElement(dstE)
	if status == notInList {
		return nil
	}
	return l.insertAfter(v, at)
}

func (l *doublyLinkedList[T]) insertBefore(v T, at *nodeElement[T]) *nodeElement[T] {
	newE := newNodeElement(v, l)
	newE.next = at
	newE.prev = at.prev
	at.prev = newE
	if newE.GetPrev() != nil {
		newE.GetPrev().(*nodeElement[T]).next = newE
	}

	l.len.Add(1)
	return newE
}

func (l *doublyLinkedList[T]) InsertBefore(v T, dstE NodeElement[T]) NodeElement[T] {
	at, status := l.checkElement(dstE)
	if status == notInList {
		return nil
	}
	return l.insertBefore(v, at)
}

func (l *doublyLinkedList[T]) Remove(targetE NodeElement[T]) NodeElement[T] {
	var (
		at     *nodeElement[T]
		status nodeElementInListStatus
	)

	// doubly linked list removes element is free to do iteration, but we have to check the value if equals.
	switch at, status = l.checkElement(targetE); status {
	case theOnlyOne:
		l.setRootHeader(l.getRoot())
		l.setRootTail(l.getRoot())
	case theFirstButNotTheLast:
		l.setRootHeader(at.next)
	case theLastButNotTheFirst:
		l.setRootTail(at.prev)
	case inMiddle:
		at.GetPrev().(*nodeElement[T]).next = at.next
		at.GetNext().(*nodeElement[T]).prev = at.prev
	default:
		return nil
	}

	at.next = nil
	at.prev = nil
	l.len.Add(-1)
	return at
}

func (l *doublyLinkedList[T]) ForEach(fn func(e NodeElement[T])) {
	if fn == nil || l.len.Load() == 0 {
		return
	}

	var iterator = l.getRoot()
	for iterator.HasNext() {
		fn(iterator.GetNext())
		iterator = iterator.GetNext()
	}
}

func (l *doublyLinkedList[T]) ReverseForEach(fn func(e NodeElement[T])) {
	if fn == nil || l.len.Load() == 0 {
		return
	}

	var iterator = l.getRoot()
	for iterator.HasPrev() {
		fn(iterator.GetPrev())
		iterator = iterator.GetPrev()
	}
}

func (l *doublyLinkedList[T]) FindFirst(targetV T, compareFn ...func(e NodeElement[T]) bool) (NodeElement[T], bool) {
	if l.len.Load() == 0 {
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

func (l *doublyLinkedList[T]) Front() NodeElement[T] {
	if l.len.Load() == 0 {
		return nil
	}
	return l.root.GetNext()
}

func (l *doublyLinkedList[T]) Back() NodeElement[T] {
	if l.len.Load() == 0 {
		return nil
	}
	return l.root.GetPrev()
}

func (l *doublyLinkedList[T]) PushFront(v T) NodeElement[T] {
	l.lazyInit()
	return l.InsertBefore(v, l.getRootHeader())
}

func (l *doublyLinkedList[T]) PushBack(v T) NodeElement[T] {
	l.lazyInit()
	return l.InsertAfter(v, l.getRootTail())
}

func (l *doublyLinkedList[T]) move(src, dst *nodeElement[T]) bool {
	if src == dst {
		return false
	}
	// Ordinarily, it is move src next to dst
	src.prev.(*nodeElement[T]).next = src.next
	src.next.(*nodeElement[T]).prev = src.prev

	src.prev = dst
	src.next = dst.next
	src.next.(*nodeElement[T]).prev = src

	dst.next = src
	return true
}

func (l *doublyLinkedList[T]) MoveToFront(targetE NodeElement[T]) {
	moved := false
	src, status := l.checkElement(targetE)
	switch status {
	case theLastButNotTheFirst:
		defer func() {
			if moved {
				l.setRootHeader(src)
				l.setRootTail(src.GetPrev())
			}
		}()
	case notInList, theOnlyOne, theFirstButNotTheLast:
		return
	default:

	}
	moved = l.move(src, &l.root)
}

func (l *doublyLinkedList[T]) MoveToBack(targetE NodeElement[T]) {
	moved := false
	src, status := l.checkElement(targetE)
	switch status {
	case theFirstButNotTheLast:
		defer func() {
			if moved {
				l.setRootHeader(src.GetNext())
				l.setRootTail(src)
			}
		}()
	case notInList, theLastButNotTheFirst, theOnlyOne:
		return
	default:

	}
	l.move(src, l.getRootTail().(*nodeElement[T]))
}

func (l *doublyLinkedList[T]) MoveBefore(srcE, dstE NodeElement[T]) {
	var (
		moved                = false
		dst, src             *nodeElement[T]
		dstStatus, srcStatus nodeElementInListStatus
	)
	switch dst, dstStatus = l.checkElement(dstE); dstStatus {
	case notInList, theOnlyOne:
		return
	default:

	}
	switch src, srcStatus = l.checkElement(srcE); srcStatus {
	case notInList, theOnlyOne:
		return
	case theFirstButNotTheLast:
		defer func() {
			if moved && dstStatus == theLastButNotTheFirst {
				l.setRootHeader(src.GetNext())
				l.setRootTail(src)
			}
			if moved && dstStatus == inMiddle {
				l.setRootHeader(src.GetNext())
			}
		}()
	case theLastButNotTheFirst:
		defer func() {
			if moved && dstStatus == theFirstButNotTheLast {
				l.setRootHeader(src)
				l.setRootTail(src.GetPrev())
			}
			if moved && dstStatus == inMiddle {
				l.setRootTail(src.GetPrev())
			}
		}()
	case inMiddle:
		defer func() {
			if moved && dstStatus == theFirstButNotTheLast {
				l.setRootHeader(src)
			}
		}()
	default:

	}
	if dst.GetPrev() == nil {
		return
	}
	dst = dst.GetPrev().(*nodeElement[T])
	moved = l.move(src, dst)
}

func (l *doublyLinkedList[T]) MoveAfter(srcE, dstE NodeElement[T]) {
	var (
		moved                = false
		dst, src             *nodeElement[T]
		dstStatus, srcStatus nodeElementInListStatus
	)
	switch dst, dstStatus = l.checkElement(dstE); dstStatus {
	case notInList, theOnlyOne:
		return
	default:

	}
	switch src, srcStatus = l.checkElement(srcE); srcStatus {
	case notInList, theOnlyOne:
		return
	case theFirstButNotTheLast:
		defer func() {
			if moved && dstStatus == theLastButNotTheFirst {
				l.setRootHeader(src.GetNext())
				l.setRootTail(src)
			}
			if moved && dstStatus == inMiddle {
				l.setRootHeader(src.GetNext())
			}
		}()
	case theLastButNotTheFirst:
		defer func() {
			if moved && dstStatus == theFirstButNotTheLast {
				l.setRootTail(src.GetPrev())
				l.setRootHeader(src)
			}
			if moved && dstStatus == inMiddle {
				l.setRootTail(src.GetPrev())
			}
		}()
	case inMiddle:
		defer func() {
			if moved && dstStatus == theLastButNotTheFirst {
				l.setRootTail(src)
			}
		}()
	default:

	}

	moved = l.move(src, dst)
}

func (l *doublyLinkedList[T]) PushFrontList(src LinkedList[T]) {
	l.lazyInit()
	for i, e := src.Len(), src.Back(); i > 0; i, e = i-1, e.GetPrev() {
		l.InsertAfter(e.GetValue(), &l.root)
	}
}

func (l *doublyLinkedList[T]) PushBackList(src LinkedList[T]) {
	l.lazyInit()
	for i, e := src.Len(), src.Front(); i > 0; i, e = i-1, e.GetNext() {
		l.Append(e)
	}
}
