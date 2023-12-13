package timer

import (
	"context"
	"errors"
	"fmt"
	"github.com/benz9527/toy-box/algo/queue"
	"log/slog"
	"runtime/debug"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	_ TimingWheel  = (*timingWheel)(nil)
	_ TimingWheels = (*xTimingWheels)(nil)
)

// 120
type timingWheel struct {
	slots []TimingWheelSlot // alignment 8, size 24; in kafka it is buckets
	// ctx is used to shut down the timing wheel and pass
	// value to control debug info.
	ctx                  context.Context                   // alignment 8, size 16
	globalDqRef          queue.DelayQueue[TimingWheelSlot] // alignment 8, size 16
	overflowWheelRef     unsafe.Pointer                    //  alignment 8, size 8; same as kafka TimingWheel(*timingWheel)
	tickMs               int64                             // alignment 8, size 8
	startMs              int64                             // alignment 8, size 8; baseline startup timestamp
	interval             int64                             // alignment 8, size 8
	currentTimeMs        int64                             // alignment 8, size 8
	slotSize             int64                             // alignment 8, size 8; in kafka it is wheelSize
	globalTaskCounterRef *atomic.Int64                     // alignment 8, size 8
	globalSlotCounterRef *atomic.Int64                     // alignment 8, size 8
}

type TimingWheelOptions func(tw *timingWheel)

func WithTimingWheelTickMs(basicTickMs time.Duration) TimingWheelOptions {
	return func(tw *timingWheel) {
		tw.tickMs = basicTickMs.Milliseconds()
	}
}

func WithTimingWheelSlotSize(slotSize int64) TimingWheelOptions {
	return func(tw *timingWheel) {
		tw.slotSize = slotSize
	}
}

func newTimingWheel(
	ctx context.Context,
	tickMs int64,
	slotSize int64,
	startMs int64,
	taskCounter *atomic.Int64,
	slotCounter *atomic.Int64,
	dq queue.DelayQueue[TimingWheelSlot],
) TimingWheel {
	tw := &timingWheel{
		ctx:                  ctx,
		tickMs:               tickMs,
		startMs:              startMs,
		slotSize:             slotSize,
		globalTaskCounterRef: taskCounter,
		globalSlotCounterRef: slotCounter,
		interval:             tickMs * slotSize,
		currentTimeMs:        startMs - (startMs % tickMs), // truncate the remainder as startMs left boundary
		slots:                make([]TimingWheelSlot, slotSize),
		globalDqRef:          dq,
	}
	// Slot initialize by doubly linked list.
	for i := int64(0); i < slotSize; i++ {
		tw.slots[i] = NewXSlot()
	}
	tw.globalSlotCounterRef.Add(slotSize)
	tw.setOverflowTimingWheel(nil)
	return tw
}

func (tw *timingWheel) GetTickMs() int64 {
	return atomic.LoadInt64(&tw.tickMs)
}

func (tw *timingWheel) GetStartMs() int64 {
	return atomic.LoadInt64(&tw.startMs)
}

func (tw *timingWheel) GetTaskCounter() int64 {
	return tw.globalTaskCounterRef.Load()
}

func (tw *timingWheel) GetCurrentTimeMs() int64 {
	return atomic.LoadInt64(&tw.currentTimeMs)
}

func (tw *timingWheel) GetInterval() int64 {
	return atomic.LoadInt64(&tw.interval)
}

func (tw *timingWheel) GetSlotSize() int64 {
	return atomic.LoadInt64(&tw.slotSize)
}

func (tw *timingWheel) getOverflowTimingWheel() TimingWheel {
	return *(*TimingWheel)(atomic.LoadPointer(&tw.overflowWheelRef))
}

func (tw *timingWheel) setOverflowTimingWheel(oftw TimingWheel) {
	atomic.StorePointer(&tw.overflowWheelRef, unsafe.Pointer(&oftw))
}

// Here related to slot level upgrade and downgrade.
func (tw *timingWheel) advanceClock(slotExpiredMs int64) {
	currentTimeMs := tw.GetCurrentTimeMs()
	tickMs := tw.GetTickMs()
	if slotExpiredMs >= currentTimeMs+tickMs {
		currentTimeMs = slotExpiredMs - (slotExpiredMs % tickMs) // truncate the remainder as slot expiredMs left boundary
		atomic.StoreInt64(&tw.currentTimeMs, currentTimeMs)      // update the current time
		oftw := tw.getOverflowTimingWheel()
		if oftw != nil {
			oftw.(*timingWheel).advanceClock(currentTimeMs)
		}
	}
}

func (tw *timingWheel) addTask(task Task, level int64) error {
	if len(task.GetJobID()) <= 0 {
		return ErrTimingWheelTaskEmptyJobID
	}
	if task.GetJob() == nil {
		return ErrTimingWheelEmptyJob
	}

	taskExpiredMs := task.GetExpiredMs()
	currentTimeMs := tw.GetCurrentTimeMs()
	tickMs := tw.GetTickMs()
	interval := tw.GetInterval()
	slotSize := tw.GetSlotSize()

	if task.Cancelled() {
		return fmt.Errorf("[timing wheel] task %s is cancelled %w",
			task.GetJobID(), ErrTimingWheelTaskCancelled)
	} else if taskExpiredMs < currentTimeMs+tickMs {
		task.setSlot(immediateExpiredSlot)
		tw.globalSlotCounterRef.Add(1)
		return fmt.Errorf("[timing wheel] task task expired ms  %d is before %d %w",
			taskExpiredMs, currentTimeMs+tickMs, ErrTimingWheelTaskIsExpired)
	} else if taskExpiredMs < currentTimeMs+interval {
		virtualID := taskExpiredMs / tickMs
		slotID := virtualID % slotSize
		slot := tw.slots[slotID]
		if !slot.setExpirationMs(virtualID * tickMs) {
			return fmt.Errorf("[timing wheel] slot (level:%d) %d unable update the expiration %w",
				level, virtualID*tickMs, ErrTimingWheelTaskUnableToBeAddedToSlot)
		}

		slot.setSlotID(slotID)
		slot.setLevel(level)
		slot.AddTask(task)
		if err := tw.globalDqRef.Offer(slot, slot.GetExpirationMs()); err != nil {
			slog.Error("[timing wheel] offer slot to delay queue error", "error", err)
		}
		tw.globalTaskCounterRef.Add(1)
		return nil
	} else {
		// Out of the interval. Put it into the higher interval timing wheel
		oftw := tw.getOverflowTimingWheel()
		if oftw == nil {
			tw.setOverflowTimingWheel(newTimingWheel(
				tw.ctx,
				interval,
				slotSize,
				currentTimeMs,
				tw.globalTaskCounterRef,
				tw.globalSlotCounterRef,
				tw.globalDqRef,
			))
		}
		// Tail recursive call, it will be free the previous stack frame.
		return tw.getOverflowTimingWheel().(*timingWheel).addTask(task, level+1)
	}
}

const (
	disableTimingWheelsSchedulePoll        = "disableTWSPoll"
	disableTimingWheelsScheduleCancelTask  = "disableTWSCancelTask"
	disableTimingWheelsScheduleExpiredSlot = "disableTWSExpSlot"
)

// size: 112
type xTimingWheels struct {
	tw           TimingWheel                       // alignment 8, size 16
	ctx          context.Context                   // alignment 8, size 16
	dq           queue.DelayQueue[TimingWheelSlot] // alignment 8, size 16; Do not use the timer.Ticker
	tasksMap     map[JobID]Task                    // alignment 8, size 8
	stopC        chan struct{}                     // alignment 8, size 8
	expiredSlotC chan TimingWheelSlot              // alignment 8, size 8
	twEventC     chan *timingWheelEvent            // alignment 8, size 8
	twEventPool  *timingWheelEventsPool            // alignment 8, size 8
	taskCounter  *atomic.Int64                     // alignment 8, size 8
	slotCounter  *atomic.Int64                     // alignment 8, size 8
	isRunning    *atomic.Bool                      // alignment 8, size 8
	// FIXME goroutine pool
}

// NewTimingWheels creates a new timing wheel.
// @param startMs the start time in milliseconds, example value time.Now().UnixMilli().
//
//	Same as the kafka, Time.SYSTEM.hiResClockMs() is used.
func NewTimingWheels(ctx context.Context, startMs int64, opts ...TimingWheelOptions) TimingWheels {
	if ctx == nil {
		return nil
	}

	xtw := &xTimingWheels{
		ctx:          ctx,
		taskCounter:  &atomic.Int64{},
		slotCounter:  &atomic.Int64{},
		twEventC:     make(chan *timingWheelEvent, 256),
		stopC:        make(chan struct{}),
		expiredSlotC: make(chan TimingWheelSlot, 3),
		tasksMap:     make(map[JobID]Task),
		isRunning:    &atomic.Bool{},
		twEventPool:  newTimingWheelEventsPool(),
	}
	xtw.isRunning.Store(false)
	tw := &timingWheel{
		startMs: startMs,
	}
	for _, o := range opts {
		if o != nil {
			o(tw)
		}
	}

	if tw.globalTaskCounterRef == nil {
		tw.globalTaskCounterRef = &atomic.Int64{}
		tw.globalTaskCounterRef.Store(0)
	}
	if tw.tickMs <= 0 {
		tw.tickMs = time.Millisecond.Milliseconds()
	}
	if tw.slotSize <= 0 {
		tw.slotSize = 20
	}
	xtw.dq = queue.NewArrayDelayQueue[TimingWheelSlot](128)
	xtw.tw = newTimingWheel(
		ctx,
		tw.tickMs,
		tw.slotSize,
		tw.startMs,
		xtw.taskCounter,
		xtw.slotCounter,
		xtw.dq,
	)
	xtw.schedule(ctx)
	return xtw
}

func (xtw *xTimingWheels) GetTickMs() int64 {
	return xtw.tw.GetTickMs()
}

func (xtw *xTimingWheels) GetStartMs() int64 {
	return xtw.tw.GetStartMs()
}

func (xtw *xTimingWheels) GetTaskCounter() int64 {
	return xtw.tw.GetTaskCounter()
}

func (xtw *xTimingWheels) GetSlotSize() int64 {
	return xtw.slotCounter.Load()
}

func (xtw *xTimingWheels) Shutdown() {
	if old := xtw.isRunning.Swap(false); !old {
		slog.Warn("[timing wheel] timing wheel is already shutdown")
		return
	}
	xtw.dq = nil
	xtw.isRunning.Store(false)

	// FIXME close on channel is no empty and will cause panic.
	close(xtw.stopC)
	close(xtw.expiredSlotC)
	close(xtw.twEventC)

	// FIXME map clear data race
}

func (xtw *xTimingWheels) AddTask(task Task) error {
	if len(task.GetJobID()) <= 0 {
		return ErrTimingWheelTaskEmptyJobID
	}
	if task.GetJob() == nil {
		return ErrTimingWheelEmptyJob
	}
	if !xtw.isRunning.Load() {
		return ErrTimingWheelStopped
	}
	event := xtw.twEventPool.Get()
	event.AddTask(task)
	xtw.twEventC <- event
	return nil
}

func (xtw *xTimingWheels) AfterFunc(delayMs time.Duration, fn Job) (Task, error) {
	if delayMs.Milliseconds() < xtw.GetTickMs() {
		return nil, fmt.Errorf("[timing wheel] delay ms %d is less than tick ms %d %w",
			delayMs.Milliseconds(), xtw.GetTickMs(), ErrTimingWheelTaskTooShortExpiration)
	}
	if fn == nil {
		return nil, ErrTimingWheelEmptyJob
	}

	now := time.Now().UTC()
	task := NewOnceTask(
		xtw.ctx,
		JobID(fmt.Sprintf("%d", now.UnixNano())), // FIXME UUID
		now.Add(delayMs).UnixMilli(),
		fn,
	)

	if !xtw.isRunning.Load() {
		return nil, ErrTimingWheelStopped
	}
	if err := xtw.AddTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (xtw *xTimingWheels) ScheduleFunc(sched Scheduler, fn Job) (Task, error) {
	if sched == nil {
		return nil, ErrTimingWheelUnknownScheduler
	}
	if fn == nil {
		return nil, ErrTimingWheelEmptyJob
	}

	now := time.Now()
	task := NewRepeatTask(
		xtw.ctx,
		JobID(fmt.Sprintf("%d", now.UnixNano())), // FIXME UUID
		now.UnixMilli(), sched,
		fn,
	)

	if !xtw.isRunning.Load() {
		return nil, ErrTimingWheelStopped
	}
	if err := xtw.AddTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (xtw *xTimingWheels) CancelTask(jobID JobID) error {
	if len(jobID) <= 0 {
		return ErrTimingWheelTaskEmptyJobID
	}

	if xtw.isRunning.Load() {
		return ErrTimingWheelStopped
	}
	task, ok := xtw.tasksMap[jobID]
	if !ok {
		return ErrTimingWheelTaskNotFound
	}

	event := xtw.twEventPool.Get()
	event.CancelTaskJobID(task.GetJobID())
	xtw.twEventC <- event
	return nil
}

func (xtw *xTimingWheels) schedule(ctx context.Context) {
	if ctx == nil {
		return
	}
	// FIXME Block error mainly caused by producer and consumer speed mismatch, lock data race.
	//  Is there any limitation mechanism could gradually  control different interval task‘s execution timeout timestamp?
	//  Tasks piling up in the same slot will cause the timing wheel to be blocked or delayed.
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("[timing wheel] event schedule panic recover", "error", err, "stack", debug.Stack())
			}
		}()
		cancelDisabled := ctx.Value(disableTimingWheelsScheduleCancelTask)
		if cancelDisabled == nil {
			cancelDisabled = false
		}
		for {
			select {
			case <-ctx.Done():
				xtw.Shutdown()
				return
			case <-xtw.stopC:
				return
			case event, ok := <-xtw.twEventC:
				if !ok {
					slog.Warn("[timing wheel] event channel has been closed")
					continue
				}
				switch event.GetOperation() {
				case addTask:
					task, ok := event.GetAddTask()
					if !ok {
						goto recycle
					}
					err := xtw.addTask(task)
					if errors.Is(err, ErrTimingWheelTaskIsExpired) {
						xtw.handleTask(task)
					}
				case cancelTask:
					jobID, ok := event.GetCancelTaskJobID()
					if !ok || cancelDisabled.(bool) {
						goto recycle
					}
					xtw.cancelTask(jobID)
				case unknown:
					fallthrough
				default:

				}
			recycle:
				xtw.twEventPool.Put(event)
			}
		}
	}()
	go func(disabled any) {
		if disabled != nil && disabled.(bool) {
			return
		}
		defer func() {
			if err := recover(); err != nil {
				slog.Error("[timing wheel] expired slot schedule panic recover", "error", err, "stack", debug.Stack())
			}
		}()
		for {
			select {
			case <-ctx.Done():
				xtw.Shutdown()
				return
			case <-xtw.stopC:
				return
			case slot, ok := <-xtw.expiredSlotC:
				if !ok {
					continue
				}
				xtw.advanceClock(slot.GetExpirationMs())
				// Here related to slot level upgrade and downgrade.
				slot.Flush(xtw.handleTask)
			}
		}
	}(ctx.Value(disableTimingWheelsScheduleExpiredSlot))
	go func(disabled any) {
		if disabled != nil && disabled.(bool) {
			return
		}
		defer func() {
			if err := recover(); err != nil {
				slog.Error("[timing wheel] poll schedule panic recover", "error", err, "stack", debug.Stack())
			}
		}()
		err := xtw.dq.PollToChannel(xtw.ctx, func() int64 {
			return time.Now().UTC().UnixMilli()
		}, xtw.expiredSlotC)
		if err != nil {
			slog.Error("[timing wheel] delay queue poll error", "error", err)
		}
		slog.Warn("[timing wheel] delay queue exit")
	}(ctx.Value(disableTimingWheelsSchedulePoll))
	xtw.isRunning.Store(true)
}

// Update all wheels' current time, in order to simulate the time is continuously incremented.
// Here related to slot level upgrade and downgrade.
func (xtw *xTimingWheels) advanceClock(timeoutMs int64) {
	xtw.tw.(*timingWheel).advanceClock(timeoutMs)
}

func (xtw *xTimingWheels) addTask(task Task) error {
	if task == nil || task.Cancelled() || !xtw.isRunning.Load() {
		return ErrTimingWheelStopped
	}
	// FIXME Recursive function to addTask a task, need to measure the performance.
	err := xtw.tw.(*timingWheel).addTask(task, 0)
	if err == nil || errors.Is(err, ErrTimingWheelTaskIsExpired) {
		// FIXME map data race
		xtw.tasksMap[task.GetJobID()] = task
	}
	return err
}

// handleTask all tasks which are called by this method
// will mean that the task must be in a slot ever and related slot
// has been expired.
func (xtw *xTimingWheels) handleTask(t Task) {
	if t == nil || !xtw.isRunning.Load() {
		return
	}

	// FIXME goroutine pool to run this.
	// [slotExpMs, slotExpMs+interval)
	var (
		prevSlotMetadata = t.GetPreviousSlotMetadata()
		slot             = t.GetSlot()
		taskLevel        int64
		runNow           bool
	)
	if prevSlotMetadata == nil {
		// Unknown task
		return
	} else {
		taskLevel = prevSlotMetadata.GetLevel()
		runNow = prevSlotMetadata.GetExpirationMs() == sentinelSlotExpiredMs || taskLevel == 0
		runNow = runNow || prevSlotMetadata.GetExpirationMs() == slotHasBeenFlushedMs &&
			t.GetExpiredMs() <= time.Now().UTC().UnixMilli()
	}
	if runNow && !t.Cancelled() {
		go t.GetJob()(xtw.ctx, t.GetJobMetadata())
	} else if t.Cancelled() {
		if slot != nil {
			slot.RemoveTask(t)
		}
		t.setSlot(nil)
		t.setSlotMetadata(nil)
		xtw.taskCounter.Add(-1)
		return
	}

	// Re-addTask loop job to timing wheel.
	// Upgrade and downgrade (move) the t from one slot to another slot.
	// Lock free.
	switch t.GetJobType() {
	case OnceJob:
		event := xtw.twEventPool.Get()
		if taskLevel == 0 && runNow {
			event.CancelTaskJobID(t.GetJobID())
			xtw.twEventC <- event
		} else if taskLevel > 0 && !runNow && xtw.isRunning.Load() == true {
			xtw.taskCounter.Add(-1)
			event.AddTask(t)
			xtw.twEventC <- event
		}
	case RepeatedJob:
		var sTask Task
		if taskLevel == 0 && runNow {
			if t.GetRestLoopCount() == 0 {
				event := xtw.twEventPool.Get()
				event.CancelTaskJobID(t.GetJobID())
				xtw.twEventC <- event
				return
			}
			_sTask, ok := t.(ScheduledTask)
			if !ok {
				return
			}
			_sTask.UpdateNextScheduledMs()
			sTask = _sTask
			if sTask.GetExpiredMs() < 0 {
				return
			}
		} else if taskLevel > 0 && !runNow {
			sTask = t
		}
		if sTask != nil && xtw.isRunning.Load() == true {
			xtw.taskCounter.Add(-1)
			event := xtw.twEventPool.Get()
			event.AddTask(sTask)
			xtw.twEventC <- event
		}
	}
	return
}

func (xtw *xTimingWheels) cancelTask(jobID JobID) {
	if !xtw.isRunning.Load() {
		return
	}

	task, ok := xtw.tasksMap[jobID]
	if !ok {
		return
	}

	if task.GetSlot() != nil && !task.GetSlot().RemoveTask(task) {
		return
	}
	task.Cancel()

	delete(xtw.tasksMap, jobID)
	xtw.taskCounter.Add(-1)
}
