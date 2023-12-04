package list

import "sync"

var (
	_ NodeElement[struct{}] = (*nodeElement[struct{}])(nil) // Type check assertion
)

// Alignment and size
// interface size is 16bytes
// pointer size is 8bytes
// string size is 16bytes
// rune size is 4bytes
// int8 size is 1byte
// int16 size is 2bytes
// int32 size is 4bytes
// int size is 8bytes
// int64 size is 8bytes
// bool size is 1byte
// byte size is 1byte
// struct{} size is 0byte

type nodeElement[T comparable] struct {
	prev, next NodeElement[T]     // alignment 8bytes * 2; size 16bytes * 2
	list       BasicLinkedList[T] // alignment 8bytes; size 16bytes
	lock       *sync.RWMutex      // alignment 8bytes; size 8bytes
	value      T                  // alignment 8bytes; size up to datatype. The type of value may be a small size type.
	// It should be placed at the end of the struct to avoid take too much padding.
}

func NewNodeElement[T comparable](v T) NodeElement[T] {
	return newNodeElement[T](v, nil)
}

func newNodeElement[T comparable](v T, list BasicLinkedList[T]) *nodeElement[T] {
	return &nodeElement[T]{
		value: v,
		list:  list,
	}
}

func NewConcurrentNodeElement[T comparable](v T) NodeElement[T] {
	return newConcurrentNodeElement[T](v, nil)
}

func newConcurrentNodeElement[T comparable](v T, list BasicLinkedList[T]) *nodeElement[T] {
	return &nodeElement[T]{
		value: v,
		list:  list,
		lock:  &sync.RWMutex{},
	}
}

func (e *nodeElement[T]) hasLock() bool {
	return e.lock != nil
}

func (e *nodeElement[T]) HasNext() bool {
	return e.next != nil
}

func (e *nodeElement[T]) HasPrev() bool {
	return e.prev != nil
}

func (e *nodeElement[T]) GetNext() NodeElement[T] {
	if e.next == nil {
		return nil
	}
	if _, ok := e.next.(*nodeElement[T]); !ok {
		return nil
	}
	return e.next
}

func (e *nodeElement[T]) GetPrev() NodeElement[T] {
	if e.prev == nil {
		return nil
	}
	if _, ok := e.prev.(*nodeElement[T]); !ok {
		return nil
	}
	return e.prev
}

func (e *nodeElement[T]) GetValue() T {
	if e.lock != nil {
		e.lock.RLock()
		defer e.lock.RUnlock()
	}
	return e.value
}

func (e *nodeElement[T]) SetValue(v T) {
	if e.lock != nil {
		e.lock.Lock()
		defer e.lock.Unlock()
	}
	e.value = v
}
