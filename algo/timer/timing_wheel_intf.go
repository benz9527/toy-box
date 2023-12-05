package timer

// Ref:
// https://dl.acm.org/doi/10.1145/41457.37504
// https://ieeexplore.ieee.org/document/650142
// github ref:
// https://github.com/RussellLuo/timingwheel

import (
	"context"
	"errors"
)

var (
	ErrTimingWheelTaskNotFound              = errors.New("task not found")
	ErrTimingWheelTaskEmptyJobID            = errors.New("empty job id in task")
	ErrTimingWheelEmptyJob                  = errors.New("empty job in task")
	ErrTimingWheelTaskIsExpired             = errors.New("task is expired")
	ErrTimingWheelTaskUnableToBeAddedToSlot = errors.New("task unable to be added to slot")
	ErrTimingWheelDuplicatedTask            = errors.New("duplicated task")
)

type TimingWheelCommonMetadata interface {
	GetTickMs() int64
	GetSlotSize() int64
	GetStartMs() int64
	GetTaskCounter() int64
}

// TimingWheel slots is private,
// they should be provided by the implementation
type TimingWheel interface {
	TimingWheelCommonMetadata
	GetInterval() int64
	GetCurrentTimeMs() int64
}

type TimingWheels interface {
	TimingWheelCommonMetadata
	AddTask(task Task) error
	CancelTask(jobID JobID) error
	Shutdown()
}

type JobID string
type Job func(ctx context.Context, jobID JobID)

// JobMetadata describes the metadata of a job
// Each slot in the timing wheel is a linked list of jobs
type JobMetadata interface {
	GetJobID() JobID
	GetExpirationMs() int64
	GetLoopCount() int64
	DecreaseLoopCount() int64
	GetJob() Job
	GetSlot() TimingWheelSlot
	SetSlot(slot TimingWheelSlot)
	Cancelled() bool
}

type Task interface {
	JobMetadata
	GetDelayMs() int64
	Cancel() bool
}

// TaskReinsert is a function that reinserts a task into the timing wheel.
// It means that the task should be executed periodically or repeatedly for a certain times.
// Reinsert will add current task to next slot, higher level slot (overflow wheel) or
// the same level slot (current wheel) depending on the expirationMs of the task.
// When the task is reinserted, the expirationMs of the task should be updated.
//  1. Check if the task is cancelled. If so, stop reinserting.
//  2. Check if the task's loop count is greater than 0. If so, decrease the loop count and reinsert.
//  3. Check if the task's loop count is -1 (run forever unless cancel manually).
//     If so, reinsert and update the expirationMs.
type TaskReinsert func(Task) // Core function

type TimingWheelSlot interface {
	GetExpirationMs() int64
	SetExpirationMs(expirationMs int64) bool
	AddTask(Task)
	RemoveTask(Task) bool
	Flush(TaskReinsert)
}
