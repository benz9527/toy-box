package timer

import (
	"context"
	"github.com/benz9527/toy-box/algo/list"
	"sync/atomic"
	"unsafe"
)

type elementTasker interface {
	getAndReleaseElementRef() list.NodeElement[Task]
}

type jobMetadata struct {
	jobID        JobID
	job          Job
	expirationMs int64
	loopCount    int64
	jobType      JobType
}

func (m *jobMetadata) GetJobID() JobID {
	return m.jobID
}

func (m *jobMetadata) GetExpiredMs() int64 {
	return m.expirationMs
}

func (m *jobMetadata) GetRestLoopCount() int64 {
	return m.loopCount
}

func (m *jobMetadata) GetJobType() JobType {
	return m.jobType
}

type task struct {
	*jobMetadata
	slotMetadata TimingWheelSlotMetadata
	slot         unsafe.Pointer // TimingWheelSlot
	// Doubly pointer reference, it is easy for us to access the element in the list.
	elementRef unsafe.Pointer // list.NodeElement[Task]
	cancelled  *atomic.Bool
}

var (
	_ Task          = (*task)(nil)
	_ elementTasker = (*task)(nil)
)

func (t *task) getAndReleaseElementRef() list.NodeElement[Task] {
	ref := t.getElementRef()
	t.setElementRef(nil)
	return ref
}

func (t *task) GetJobID() JobID {
	return t.jobID
}

func (t *task) GetJobMetadata() JobMetadata {
	md := &jobMetadata{
		jobID:        t.jobID,
		job:          t.job,
		expirationMs: t.expirationMs,
		loopCount:    t.loopCount,
		jobType:      t.jobType,
	}
	return md
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

func (t *task) GetExpiredMs() int64 {
	return atomic.LoadInt64(&t.expirationMs)
}

func (t *task) GetSlot() TimingWheelSlot {
	return *(*TimingWheelSlot)(atomic.LoadPointer(&t.slot))
}

func (t *task) setSlot(slot TimingWheelSlot) {
	atomic.StorePointer(&t.slot, unsafe.Pointer(&slot))
}

func (t *task) GetPreviousSlotMetadata() TimingWheelSlotMetadata {
	return t.slotMetadata
}

func (t *task) setSlotMetadata(slotMetadata TimingWheelSlotMetadata) {
	t.slotMetadata = slotMetadata
}

func (t *task) getElementRef() list.NodeElement[Task] {
	return *(*list.NodeElement[Task])(atomic.LoadPointer(&t.elementRef))
}

func (t *task) setElementRef(elementRef list.NodeElement[Task]) {
	atomic.StorePointer(&t.elementRef, unsafe.Pointer(&elementRef))
}

func (t *task) GetJobType() JobType {
	return t.jobType
}

func (t *task) Cancel() bool {
	if stopped := t.cancelled.Swap(true); stopped {
		// Previous value is true, it means that the task has been cancelled.
		return true
	}

	// If task is cancelled, it will be removed from the timing wheel automatically
	// in other process, so we don't need to remove it here.

	if t.jobType == OnceJob {
		atomic.SwapInt64(&t.loopCount, 0)
	} else if t.jobType == RepeatedJob {
		atomic.SwapInt64(&t.loopCount, t.GetRestLoopCount()-1)
	}
	return true
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
	atomic.SwapInt64(&t.beginMs, expiredMs)
}

func (t *xScheduledTask) GetRestLoopCount() int64 {
	return t.scheduler.GetRestLoopCount()
}

func (t *xScheduledTask) getAndReleaseElementRef() list.NodeElement[Task] {
	return t.task.getAndReleaseElementRef()
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
			jobMetadata: &jobMetadata{
				jobID:        jobID,
				expirationMs: expiredMs,
				loopCount:    1,
				job:          job,
				jobType:      OnceJob,
			},
			cancelled: &atomic.Bool{},
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
				jobMetadata: &jobMetadata{
					jobID:   jobID,
					job:     job,
					jobType: RepeatedJob,
				},
				cancelled: &atomic.Bool{},
			},
			ctx: ctx,
		},
		scheduler: scheduler,
		beginMs:   beginMs,
	}
	t.UpdateNextScheduledMs()
	return t
}
