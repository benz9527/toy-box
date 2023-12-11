package timer

import (
	"github.com/benz9527/toy-box/algo/list"
	"sync/atomic"
)

type slotMetadata struct {
	expirationMs int64
	slotID       int64
	level        int64
}

func (slot *slotMetadata) GetExpirationMs() int64 {
	return atomic.LoadInt64(&slot.expirationMs)
}

func (slot *slotMetadata) setExpirationMs(expirationMs int64) bool {
	// If expirationMs is -1, it means that the slot is empty and will be
	// reused by the timing wheel.
	return atomic.SwapInt64(&slot.expirationMs, expirationMs) != expirationMs
}

func (slot *slotMetadata) GetSlotID() int64 {
	return atomic.LoadInt64(&slot.slotID)
}

func (slot *slotMetadata) setSlotID(slotID int64) {
	atomic.SwapInt64(&slot.slotID, slotID)
}

func (slot *slotMetadata) GetLevel() int64 {
	return atomic.LoadInt64(&slot.level)
}

func (slot *slotMetadata) setLevel(level int64) {
	atomic.SwapInt64(&slot.level, level)
}

// xSlot the segment rwlock free
type xSlot struct {
	*slotMetadata
	tasks list.LinkedList[Task]
}

const (
	slotHasBeenFlushedMs  = -1
	sentinelSlotExpiredMs = -2
)

var (
	// immediateExpiredSlot is a sentinel slot that is used to mark the tasks should be executed immediately.
	// There are without any tasks in this slot actual.
	immediateExpiredSlot = newSentinelSlot()
)

func NewXSlot() TimingWheelSlot {
	return &xSlot{
		slotMetadata: &slotMetadata{
			expirationMs: -1,
		},
		tasks: list.NewLinkedList[Task](),
	}
}

func newSentinelSlot() TimingWheelSlot {
	return &xSlot{
		slotMetadata: &slotMetadata{
			expirationMs: sentinelSlotExpiredMs,
		},
	}
}

func (slot *xSlot) GetMetadata() TimingWheelSlotMetadata {
	metadata := *slot.slotMetadata // Copy instead of reference
	return &metadata
}

func (slot *xSlot) AddTask(task Task) {
	if task == nil && slot.GetExpirationMs() == slotHasBeenFlushedMs {
		return
	}

	elementRefs := slot.tasks.AppendValue(task)
	task.setSlot(slot)
	task.setSlotMetadata(slot.GetMetadata())
	switch _task := task.(type) {
	case *xScheduledTask:
		_task.setElementRef(elementRefs[0])
	case *xTask:
		_task.setElementRef(elementRefs[0])
	}
}

// removeTask clear the reference of the task, avoid memory leak.
// Reserve the previous slot metadata of the task, it is used to
// do reinsert or other operations, avoid data race and nil panic.
func (slot *xSlot) removeTask(task Task) bool {
	if slot == immediateExpiredSlot || task == nil || task.GetSlot() != slot {
		return false
	}
	// Remove task from slot but not cancel it, lock free
	e := slot.tasks.Remove(task.(elementTasker).getAndReleaseElementRef())
	if e == nil {
		return false
	}
	task.setSlot(nil) // clean reference, avoid memory leak
	return true
}

// RemoveTask the slot must be not expired.
func (slot *xSlot) RemoveTask(task Task) bool {
	if slot.GetExpirationMs() == slotHasBeenFlushedMs {
		return false
	}
	task.setSlotMetadata(nil)
	return slot.removeTask(task)
}

// Flush Timing wheel scheduling algorithm core function
func (slot *xSlot) Flush(reinsert TaskHandler) {
	// Reset the slot, ready for next round.
	slot.setExpirationMs(slotHasBeenFlushedMs)
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
		slot.removeTask(task) // clean task reference at first
		reinsert(task)
	})
}
