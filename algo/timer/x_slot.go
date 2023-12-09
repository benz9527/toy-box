package timer

import (
	"github.com/benz9527/toy-box/algo/list"
	"sync"
	"sync/atomic"
)

type xSlot struct {
	expirationMs int64
	slotID       int64
	level        int64
	tasks        list.LinkedList[Task]
	lock         *sync.RWMutex
}

const (
	sentinelSlotExpiredMs = -2
)

var (
	immediateExpiredSlot = newSentinelSlot()
)

func NewXSlot() TimingWheelSlot {
	return &xSlot{
		expirationMs: -1,
		tasks:        list.NewLinkedList[Task](),
		lock:         &sync.RWMutex{},
	}
}

func newSentinelSlot() TimingWheelSlot {
	return &xSlot{
		expirationMs: sentinelSlotExpiredMs,
	}
}

func (slot *xSlot) GetExpirationMs() int64 {
	return atomic.LoadInt64(&slot.expirationMs)
}

func (slot *xSlot) setExpirationMs(expirationMs int64) bool {
	// If expirationMs is -1, it means that the slot is empty and will be
	// reused by the timing wheel.
	return atomic.SwapInt64(&slot.expirationMs, expirationMs) != expirationMs
}

func (slot *xSlot) GetSlotID() int64 {
	return atomic.LoadInt64(&slot.slotID)
}

func (slot *xSlot) setSlotID(slotID int64) {
	atomic.SwapInt64(&slot.slotID, slotID)
}

func (slot *xSlot) GetLevel() int64 {
	return atomic.LoadInt64(&slot.level)
}

func (slot *xSlot) setLevel(level int64) {
	atomic.SwapInt64(&slot.level, level)
}

func (slot *xSlot) AddTask(task Task) {
	if task == nil {
		return
	}

	slot.lock.Lock()
	elementRefs := slot.tasks.AppendValue(task)
	slot.lock.Unlock()

	task.setSlot(slot)
	switch _task := task.(type) {
	case *xScheduledTask:
		_task.setElementRef(elementRefs[0])
	case *xTask:
		_task.setElementRef(elementRefs[0])
	}
}

func (slot *xSlot) removeTask(task Task) bool {
	if slot == immediateExpiredSlot || task == nil || task.GetSlot() != slot {
		return false
	}
	// Remove task from slot but not cancel it
	slot.lock.Lock()
	defer slot.lock.Unlock()
	e := slot.tasks.Remove(task.(elementTasker).getAndReleaseElementRef())
	if e == nil {
		return false
	}
	task.setSlot(nil) // clean reference, avoid memory leak
	return true
}

func (slot *xSlot) RemoveTask(task Task) bool {
	return slot.removeTask(task)
}

// Flush Timing wheel scheduling algorithm core function
func (slot *xSlot) Flush(reinsert TaskReinsert) {
	// Due to the slot has been expired, we have to handle all tasks from the slot.
	// 1. If the task is cancelled, we will remove it from the slot.
	// 2. If the task is not cancelled:
	//  2.1 Check the task is a high level timing wheel task or not.
	//      If so, reinsert the task to the lower level timing wheel.
	//      Otherwise, run the task.
	//  2.2 If the task is a low level timing wheel task, run the task.
	//      If the task is a repeat task, reinsert the task to the current timing wheel.
	//      Otherwise, cancel it.
	// 3. Remove the tasks from the slot.
	// 4. Reset the slot, ready for next round.
	slot.tasks.ForEach(func(idx int64, iterator list.NodeElement[Task]) {
		task := iterator.GetValue()
		reinsert(task)
		slot.removeTask(task)
	})
	// Reset the slot, ready for next round.
	slot.setExpirationMs(-1)
}
