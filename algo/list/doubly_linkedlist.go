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
	emtpyList
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

func (l *doublyLinkedList[T]) getRootHead() NodeElement[T] {
	return l.root.next
}

func (l *doublyLinkedList[T]) setRootHead(targetE NodeElement[T]) {
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
	l.setRootHead(&l.root)
	l.setRootTail(&l.root)
	l.len.Store(0)
	return l
}

func (l *doublyLinkedList[T]) Len() int64 {
	return l.len.Load()
}

func (l *doublyLinkedList[T]) checkElement(targetE NodeElement[T]) (*nodeElement[T], nodeElementInListStatus) {
	if targetE == nil {
		return nil, notInList
	}
	at, ok := targetE.(*nodeElement[T])
	if !ok || at.list != l || at.hasLock() != l.root.hasLock() {
		return nil, notInList
	}
	if l.len.Load() == 0 {
		return l.getRoot().(*nodeElement[T]), emtpyList
	}

	// mem address compare
	switch {
	case at.GetPrev() == nil && at.GetNext() == nil:
		// targetE is the first one and the last one
		if l.getRootHead() != targetE || l.getRootTail() != targetE {
			return nil, notInList
		}
		return at, theOnlyOne
	case at.GetPrev() == nil && at.GetNext() != nil:
		// targetE is the first one but not the last one
		if at.GetNext().GetPrev() != targetE {
			return nil, notInList
		}
		// Ignore l.setRootTail (tail)
		return at, theFirstButNotTheLast
	case at.GetPrev() != nil && at.GetNext() == nil:
		// targetE is the last one but not the first one
		if at.GetPrev().GetNext() != targetE {
			return nil, notInList
		}
		// Ignore l.setRootHead (head)
		return at, theLastButNotTheFirst
	case at.GetPrev() != nil && at.GetNext() != nil:
		// targetE is neither the first one nor the last one
		if at.GetPrev().GetNext() != targetE && at.GetNext().GetPrev() != targetE {
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
		l.setRootHead(e)
		l.setRootTail(e)
		l.len.Add(1)
		return e
	}

	lastOne := l.getRootTail().(*nodeElement[T])
	lastOne.next = e
	e.prev, e.next = lastOne, nil
	l.setRootTail(e)

	l.len.Add(1)
	return e
}

func (l *doublyLinkedList[T]) Append(elements ...NodeElement[T]) []NodeElement[T] {
	if l.getRootTail() == nil {
		return nil
	}
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
	// FIXME How to decrease the memory allocation for each operation?
	//  sync.Pool reused the released objects?
	if len(values) <= 0 {
		return nil
	} else if len(values) == 1 {
		return l.Append(newNodeElement(values[0], l))
	}

	newElements := make([]NodeElement[T], 0, len(values))
	for _, v := range values {
		newElements = append(newElements, newNodeElement(v, l))
	}
	return l.Append(newElements...)
}

func (l *doublyLinkedList[T]) insertAfter(newE, at *nodeElement[T]) *nodeElement[T] {
	if at == l.getRoot() {
		newE.prev = nil
		newE.next = nil
		l.setRootHead(newE)
		l.setRootTail(newE)
	} else {
		newE.prev = at
		newE.next = at.GetNext()
		at.next = newE
		if newE.GetNext() != nil {
			newE.GetNext().(*nodeElement[T]).prev = newE
		}
	}
	l.len.Add(1)
	return newE
}

func (l *doublyLinkedList[T]) InsertAfter(v T, dstE NodeElement[T]) NodeElement[T] {
	at, status := l.checkElement(dstE)
	if status == notInList {
		return nil
	}

	newE := newNodeElement(v, l)
	newE = l.insertAfter(newE, at)
	switch status {
	case theOnlyOne, theLastButNotTheFirst:
		l.setRootTail(newE)
	default:

	}
	return newE
}

func (l *doublyLinkedList[T]) insertBefore(newE, at *nodeElement[T]) *nodeElement[T] {
	if at == l.getRoot() {
		newE.prev = nil
		newE.next = nil
		l.setRootHead(newE)
		l.setRootTail(newE)
	} else {
		newE.next = at
		newE.prev = at.prev
		at.prev = newE
		if newE.GetPrev() != nil {
			newE.GetPrev().(*nodeElement[T]).next = newE
		}
	}
	l.len.Add(1)
	return newE
}

func (l *doublyLinkedList[T]) InsertBefore(v T, dstE NodeElement[T]) NodeElement[T] {
	at, status := l.checkElement(dstE)
	if status == notInList {
		return nil
	}
	newE := newNodeElement(v, l)
	newE = l.insertBefore(newE, at)
	switch status {
	case theOnlyOne, theFirstButNotTheLast:
		l.setRootHead(newE)
	default:

	}
	return newE
}

func (l *doublyLinkedList[T]) Remove(targetE NodeElement[T]) NodeElement[T] {
	var (
		at     *nodeElement[T]
		status nodeElementInListStatus
	)

	// doubly linked list removes element is free to do iteration, but we have to check the value if equals.
	// avoid memory leaks
	switch at, status = l.checkElement(targetE); status {
	case theOnlyOne:
		l.setRootHead(l.getRoot())
		l.setRootTail(l.getRoot())
	case theFirstButNotTheLast:
		l.setRootHead(at.GetNext())
		at.GetNext().(*nodeElement[T]).prev = nil
	case theLastButNotTheFirst:
		l.setRootTail(at.GetPrev())
		at.GetPrev().(*nodeElement[T]).next = nil
	case inMiddle:
		at.GetPrev().(*nodeElement[T]).next = at.next
		at.GetNext().(*nodeElement[T]).prev = at.prev
	default:
		return nil
	}

	// avoid memory leaks
	at.list = nil
	at.next = nil
	at.prev = nil
	if at.lock != nil {
		at.lock = nil
	}
	l.len.Add(-1)
	return at
}

func (l *doublyLinkedList[T]) ForEach(fn func(idx int64, e NodeElement[T])) {
	if fn == nil || l.len.Load() == 0 ||
		l.getRoot() == l.getRootHead() && l.getRoot() == l.getRootTail() {
		return
	}

	var (
		iterator       = l.getRoot().GetNext()
		idx      int64 = 0
	)
	for iterator != nil {
		// Avoid remove in iteration result in memory leak
		next := iterator.GetNext()
		fn(idx, iterator)
		iterator, idx = next, idx+1
	}
}

func (l *doublyLinkedList[T]) ReverseForEach(fn func(idx int64, e NodeElement[T])) {
	if fn == nil || l.len.Load() == 0 ||
		l.getRoot() == l.getRootHead() && l.getRoot() == l.getRootTail() {
		return
	}

	var (
		iterator       = l.getRoot().GetPrev()
		idx      int64 = 0
	)
	for iterator != nil {
		// Avoid remove in iteration result in memory leak
		prev := iterator.GetPrev()
		fn(idx, iterator)
		iterator, idx = prev, idx+1
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
	return l.InsertBefore(v, l.getRootHead())
}

func (l *doublyLinkedList[T]) PushBack(v T) NodeElement[T] {
	return l.InsertAfter(v, l.getRootTail())
}

func (l *doublyLinkedList[T]) move(src, dst *nodeElement[T]) bool {
	if src == dst {
		return false
	}
	// Ordinarily, it is move src next to dst
	if src.HasPrev() { // src is not the first node
		src.GetPrev().(*nodeElement[T]).next = src.next
	}
	if src.HasNext() { // src is not the last node
		src.GetNext().(*nodeElement[T]).prev = src.prev
	}

	if dst == l.getRoot() { // dst is the first node
		src.prev = nil
	} else {
		src.prev = dst
	}
	src.next = dst.GetNext()
	if src.HasNext() { // src updated but not the last node
		src.GetNext().(*nodeElement[T]).prev = src
	}

	dst.next = src
	return true
}

func (l *doublyLinkedList[T]) MoveToFront(targetE NodeElement[T]) {
	moved := false
	src, status := l.checkElement(targetE)
	switch status {
	case theLastButNotTheFirst:
		defer func(prev NodeElement[T]) {
			if moved {
				l.setRootHead(src)
				l.setRootTail(prev)
			}
		}(src.GetPrev()) // register the element immediately
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
		defer func(next NodeElement[T]) {
			if moved {
				l.setRootHead(next)
				l.setRootTail(src)
			}
		}(src.GetNext()) // register the element immediately
	case notInList, theLastButNotTheFirst, theOnlyOne:
		return
	default:

	}
	moved = l.move(src, l.getRootTail().(*nodeElement[T]))
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
		defer func(next NodeElement[T]) {
			// src from the head to the middle
			if moved && dstStatus == theLastButNotTheFirst || dstStatus == inMiddle {
				l.setRootHead(next)
			}
		}(src.GetNext()) // register the element immediately
	case theLastButNotTheFirst:
		defer func(prev NodeElement[T]) {
			// src from the tail to the head
			if moved && dstStatus == theFirstButNotTheLast {
				l.setRootHead(src)
				l.setRootTail(prev)
			}
			if moved && dstStatus == inMiddle {
				l.setRootTail(prev)
			}
		}(src.GetPrev()) // register the element immediately
	case inMiddle:
		defer func() {
			// src as the new first node
			if moved && dstStatus == theFirstButNotTheLast {
				l.setRootHead(src)
			}
		}()
	default:

	}
	if dst.GetPrev() == nil {
		return
	}
	dstPrev := dst.GetPrev().(*nodeElement[T])
	moved = l.move(src, dstPrev)
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
		defer func(next NodeElement[T]) {
			// src from the head to the tail
			if moved && dstStatus == theLastButNotTheFirst {
				l.setRootHead(next)
				l.setRootTail(src)
			}
			// src from the head to the middle
			if moved && dstStatus == inMiddle {
				l.setRootHead(next)
			}
		}(src.GetNext()) // register the element immediately
	case theLastButNotTheFirst:
		defer func(prev NodeElement[T]) {
			// src from the tail to the middle
			if moved && dstStatus == theFirstButNotTheLast || dstStatus == inMiddle {
				l.setRootTail(prev)
			}
		}(src.GetPrev())
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
	if dl, ok := src.(*doublyLinkedList[T]); !ok || dl != nil && dl.getRoot() == l.getRoot() {
		// avoid type mismatch and self copy
		return
	}

	for i, e := src.Len(), src.Back(); i > 0; i-- {
		prev := e.GetPrev()
		// Clean the node element reference, not a required operation
		e.(*nodeElement[T]).list = l
		e.(*nodeElement[T]).prev = nil
		e.(*nodeElement[T]).next = nil
		l.insertAfter(e.(*nodeElement[T]), &l.root)
		e = prev
	}
}

func (l *doublyLinkedList[T]) PushBackList(src LinkedList[T]) {
	if dl, ok := src.(*doublyLinkedList[T]); !ok || dl != nil && dl.getRoot() == l.getRoot() {
		// avoid type mismatch and self copy
		return
	}
	for i, e := src.Len(), src.Front(); i > 0; i-- {
		next := e.GetNext()
		// Clean the node element reference, not a required operation
		e.(*nodeElement[T]).list = l
		e.(*nodeElement[T]).prev = nil
		e.(*nodeElement[T]).next = nil
		l.Append(e)
		e = next
	}
}
