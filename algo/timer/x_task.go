package timer

import (
	"context"
	"github.com/benz9527/toy-box/algo/list"
	"sync"
	"sync/atomic"
)

type elementTasker interface {
	getAndReleaseElementRef() list.NodeElement[Task]
}

type task struct {
	jobID        JobID
	job          Job
	expirationMs int64
	slot         TimingWheelSlot
	// Doubly pointer reference, it is easy for us to access the element in the list.
	elementRef list.NodeElement[Task]

	loopCount int64
	cancelled atomic.Bool
	taskType  TaskType
	lock      *sync.RWMutex
}

var (
	_ Task          = (*task)(nil)
	_ elementTasker = (*task)(nil)
)

func (t *task) getAndReleaseElementRef() list.NodeElement[Task] {
	t.lock.Lock()
	defer t.lock.Unlock()
	ref := t.elementRef
	t.elementRef = nil
	return ref
}

func (t *task) GetJobID() JobID {
	return t.jobID
}

func (t *task) GetRestLoopCount() int64 {
	return atomic.LoadInt64(&t.loopCount)
}

func (t *task) GetJob() Job {
	return t.job
}

func (t *task) Cancelled() bool {
	return t.cancelled.Load()
}

func (t *task) GetExpirationMs() int64 {
	return atomic.LoadInt64(&t.expirationMs)
}

func (t *task) GetSlot() TimingWheelSlot {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.slot
}

func (t *task) setSlot(slot TimingWheelSlot) {
	if slot == nil {
		return
	}
	t.lock.Lock()
	t.slot = slot
	t.lock.Unlock()
}

func (t *task) getElementRef() list.NodeElement[Task] {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.elementRef
}

func (t *task) setElementRef(elementRef list.NodeElement[Task]) {
	t.lock.Lock()
	t.elementRef = elementRef
	t.lock.Unlock()
}

func (t *task) GetTaskType() TaskType {
	return t.taskType
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
		if t.cancelled.Swap(true) {
			// Previous value is true, it means that the task has been cancelled.
			stopped = true
			break
		}
	}
	if stopped && t.taskType == OnceTask {
		atomic.SwapInt64(&t.loopCount, 0)
	} else if stopped && t.taskType == RepeatTask {
		atomic.SwapInt64(&t.loopCount, t.GetRestLoopCount()-1)
	}
	return stopped
}

type xTask struct {
	*task
	ctx context.Context
}

func (t *xTask) getAndReleaseElementRef() list.NodeElement[Task] {
	return t.task.getAndReleaseElementRef()
}

func (t *xTask) Cancel() bool {
	return t.task.Cancel()
}

type xScheduledTask struct {
	*xTask
	beginMs   int64
	scheduler Scheduler
}

func (t *xScheduledTask) UpdateNextScheduledMs() {
	expiredMs := t.scheduler.next(t.beginMs)
	atomic.StoreInt64(&t.expirationMs, expiredMs)
	if expiredMs == -1 {
		return
	}
	t.beginMs = expiredMs
}

func (t *xScheduledTask) GetRestLoopCount() int64 {
	return t.scheduler.GetRestLoopCount()
}

func (t *xScheduledTask) getAndReleaseElementRef() list.NodeElement[Task] {
	return t.elementRef
}

func (t *xScheduledTask) Cancel() bool {
	return t.task.Cancel()
}

var (
	_ Task          = (*task)(nil)
	_ Task          = (*xTask)(nil)
	_ ScheduledTask = (*xScheduledTask)(nil)
)

func NewOnceTask(
	ctx context.Context,
	jobID JobID,
	expiredMs int64,
	job Job,
) Task {
	if ctx == nil {
		return nil
	}

	t := &xTask{
		task: &task{
			jobID:        jobID,
			expirationMs: expiredMs,
			loopCount:    1,
			job:          job,
			lock:         &sync.RWMutex{},
		},
		ctx: ctx,
	}
	return t
}

func NewRepeatTask(
	ctx context.Context,
	jobID JobID,
	beginMs int64,
	scheduler Scheduler,
	job Job,
) ScheduledTask {
	if ctx == nil || scheduler == nil || job == nil {
		return nil
	}
	t := &xScheduledTask{
		xTask: &xTask{
			task: &task{
				jobID: jobID,
				job:   job,
				lock:  &sync.RWMutex{},
			},
			ctx: ctx,
		},
		scheduler: scheduler,
		beginMs:   beginMs,
	}
	t.UpdateNextScheduledMs()
	return t
}
