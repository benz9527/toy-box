package timer

import (
	"github.com/benz9527/toy-box/algo/list"
	"sync"
	"sync/atomic"
)

type xSlot struct {
	expirationMs int64
	tasks        list.LinkedList[Task]
	lock         *sync.RWMutex
}

func NewXSlot() TimingWheelSlot {
	return &xSlot{
		expirationMs: -1,
		tasks:        list.NewLinkedList[Task](),
		lock:         &sync.RWMutex{},
	}
}

func (slot *xSlot) GetExpirationMs() int64 {
	return atomic.LoadInt64(&slot.expirationMs)
}

func (slot *xSlot) SetExpirationMs(expirationMs int64) bool {
	// If expirationMs is -1, it means that the slot is empty and will be
	// reused by the timing wheel.
	return atomic.SwapInt64(&slot.expirationMs, expirationMs) != expirationMs
}

func (slot *xSlot) AddTask(task Task) {
	if task == nil {
		return
	}
	if _, ok := task.(*xTask); !ok {
		return
	}

	slot.lock.Lock()
	defer slot.lock.Unlock()
	elementRef := slot.tasks.PushBack(task)
	task.SetSlot(slot)
	task.(*xTask).elementRef = elementRef
}

func (slot *xSlot) removeTask(task Task) bool {
	if task == nil {
		return false
	}
	if task.GetSlot() != slot {
		return false
	}
	slot.tasks.Remove(task.(*xTask).elementRef)
	// clean reference, avoid memory leak
	task.SetSlot(nil)
	task.(*xTask).elementRef = nil
	return true
}

func (slot *xSlot) RemoveTask(task Task) bool {
	slot.lock.Lock()
	defer slot.lock.Unlock()
	return slot.removeTask(task)
}

func (slot *xSlot) Flush(reinsert TaskReinsert) {
	slot.lock.Lock()
	defer slot.lock.Unlock()

	// Fetch all tasks in the slot.
	// Because the tasks are all expired.
	iterator := slot.tasks.Front()
	for iterator != nil && iterator.HasNext() {
		next := iterator.GetNext()
		task := iterator.GetValue().(*task)
		slot.removeTask(task)
		reinsert(task)
		iterator = next
	}
	slot.SetExpirationMs(-1)
}
