package timer

import (
	"context"
	"github.com/benz9527/toy-box/algo/list"
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"
)

type task struct {
	jobID        JobID
	job          Job
	expirationMs int64
	// TimingWheelSlot atomic store and load is safe for concurrent access.
	slot unsafe.Pointer
	// Doubly pointer reference, it is easy for us to access the element in the list.
	elementRef list.NodeElement[Task]

	loopCount int64
	cancelled atomic.Bool
}

func (t *task) GetJobID() JobID {
	return t.jobID
}

func (t *task) GetLoopCount() int64 {
	return atomic.LoadInt64(&t.loopCount)
}

func (t *task) IncreaseLoopCount() int64 {
	return atomic.SwapInt64(&t.loopCount, t.GetLoopCount()+1)
}

func (t *task) DecreaseLoopCount() int64 {
	return atomic.SwapInt64(&t.loopCount, t.GetLoopCount()-1)
}

func (t *task) GetJob() Job {
	return t.job
}

func (t *task) Cancelled() bool {
	return t.cancelled.Load()
}

func (t *task) GetDelayMs() int64 {
	return 0
}

func (t *task) GetExpirationMs() int64 {
	return atomic.LoadInt64(&t.expirationMs)
}

func (t *task) SetExpirationMs(expirationMs int64) {
	t.expirationMs = expirationMs
}

func (t *task) GetSlot() TimingWheelSlot {
	ptr := (*any)(atomic.LoadPointer(&t.slot))
	return (*ptr).(TimingWheelSlot)
}

func (t *task) SetSlot(slot TimingWheelSlot) {
	v := reflect.ValueOf(slot)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	atomic.StorePointer(&t.slot, v.UnsafePointer())
}

func (t *task) Cancel() bool {
	stopped := false
	for slot := t.GetSlot(); slot != nil && !t.Cancelled(); slot = t.GetSlot() {
		// Remove t from all timing wheel.
		// Actually, there is only one slot for the t.
		// In this loop, we avoid 2 scenarios:
		// 1. We remove the t from current slot by below method.
		//  But at the same time, the t is expired and reinserted to the next slot,
		//  which handled by slot.Flush() in other goroutine.
		stopped = slot.RemoveTask(t)
		if old := t.cancelled.Swap(true); old {
			stopped = true
			break
		}
	}
	return stopped
}

type xTask struct {
	task
	ctx context.Context
}

var (
	_ Task = (*task)(nil)
	_ Task = (*xTask)(nil)
)

func NewXTask(
	ctx context.Context,
	jobID JobID,
	expiredMs time.Duration,
	loopCount int64,
	job Job,
) Task {
	if ctx == nil {
		return nil
	}

	t := &xTask{
		task: task{
			jobID:        jobID,
			expirationMs: expiredMs.Milliseconds(),
			loopCount:    loopCount,
			job:          job,
		},
		ctx: ctx,
	}
	return t
}
