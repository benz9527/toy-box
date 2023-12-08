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
	atomic.StoreInt64(&slot.slotID, slotID)
}

func (slot *xSlot) GetLevel() int64 {
	return atomic.LoadInt64(&slot.level)
}

func (slot *xSlot) setLevel(level int64) {
	atomic.StoreInt64(&slot.level, level)
}

func (slot *xSlot) AddTask(task Task) {
	if task == nil {
		return
	}
	if _, ok := task.(*xTask); !ok {
		return
	}

	slot.lock.Lock()
	elementRefs := slot.tasks.AppendValue(task)
	slot.lock.Unlock()

	task.setSlot(slot)
	task.(*xTask).setElementRef(elementRefs[0])
}

func (slot *xSlot) removeTask(task Task) bool {
	if slot == immediateExpiredSlot || task == nil {
		return false
	}
	if task.GetSlot() != slot {
		return false
	}
	// remove task from slot but not cancel it
	slot.lock.Lock()
	slot.tasks.Remove(task.(elementTasker).getAndReleaseElementRef())
	slot.lock.Unlock()
	// clean reference, avoid memory leak
	task.setSlot(nil)
	return true
}

func (slot *xSlot) RemoveTask(task Task) bool {
	if task == nil || task.Cancelled() {
		return true
	}
	return slot.removeTask(task)
}

func (slot *xSlot) Flush(reinsert TaskReinsert) {
	//slot.lock.Lock()
	//defer slot.lock.Unlock()

	// Fetch all tasks in the slot.
	// Because the tasks are all expired.
	slot.tasks.ForEach(func(idx int64, iterator list.NodeElement[Task]) {
		// FIXME lock race
		task := iterator.GetValue()
		reinsert(task)
		slot.removeTask(task)
	})
	// Reset the slot, ready for next round.
	slot.setExpirationMs(-1)
}
