package timer

import (
	"context"
	"errors"
	"fmt"
	"github.com/benz9527/toy-box/algo/queue"
	"log/slog"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_ TimingWheel  = (*timingWheel)(nil)
	_ TimingWheels = (*xTimingWheels)(nil)
)

// 136
type timingWheel struct {
	slots []TimingWheelSlot // alignment 8, size 24; in kafka it is buckets
	// ctx is used to shut down the timing wheel and pass
	// value to control debug info.
	ctx                  context.Context                   // alignment 8, size 16
	overflowWheelRef     TimingWheel                       // alignment 8, size 16; same as kafka
	globalDqRef          queue.DelayQueue[TimingWheelSlot] // alignment 8, size 16
	tickMs               int64                             // alignment 8, size 8
	startMs              int64                             // alignment 8, size 8; baseline startup timestamp
	interval             int64                             // alignment 8, size 8
	currentTimeMs        int64                             // alignment 8, size 8
	slotSize             int64                             // alignment 8, size 8; in kafka it is wheelSize
	globalTaskCounterRef *atomic.Int64                     // alignment 8, size 8
	globalSlotCounterRef *atomic.Int64                     // alignment 8, size 8
	lock                 *sync.RWMutex                     // alignment 8, size 8
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
		lock:                 &sync.RWMutex{},
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

// Here related to slot level upgrade and downgrade.
func (tw *timingWheel) advanceClock(slotExpiredMs int64) {
	currentTimeMs := tw.GetCurrentTimeMs()
	tickMs := tw.GetTickMs()
	if slotExpiredMs >= currentTimeMs+tickMs {
		currentTimeMs = slotExpiredMs - (slotExpiredMs % tickMs) // truncate the remainder as slot expiredMs left boundary
		atomic.StoreInt64(&tw.currentTimeMs, currentTimeMs)      // update the current time
		if tw.overflowWheelRef != nil {
			if otw, ok := tw.overflowWheelRef.(*timingWheel); ok {
				otw.advanceClock(currentTimeMs)
			}
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

		tw.lock.Lock()
		slot.AddTask(task)
		slot.setSlotID(slotID)
		slot.setLevel(level)
		if err := tw.globalDqRef.Offer(slot, slot.GetExpirationMs()); err != nil {
			slog.Error("[timing wheel] offer slot to delay queue error", "error", err)
		}
		tw.lock.Unlock()
		tw.globalTaskCounterRef.Add(1)
		return nil
	} else {
		// Out of the interval. Put it into the higher interval timing wheel
		if tw.overflowWheelRef == nil {
			tw.lock.Lock()
			tw.overflowWheelRef = newTimingWheel(
				tw.ctx,
				interval,
				slotSize,
				currentTimeMs,
				tw.globalTaskCounterRef,
				tw.globalSlotCounterRef,
				tw.globalDqRef,
			)
			tw.lock.Unlock()
		}
		// Tail recursive call, it will be free the previous stack frame.
		return tw.overflowWheelRef.(*timingWheel).addTask(task, level+1)
	}
}

const (
	disableTimingWheelsSchedulePoll        = "disableTWSPoll"
	disableTimingWheelsScheduleCancelTask  = "disableTWSCancelTask"
	disableTimingWheelsScheduleExpiredSlot = "disableTWSExpSlot"
)

// size: 120
type xTimingWheels struct {
	tw           TimingWheel                       // alignment 8, size 16
	ctx          context.Context                   // alignment 8, size 16
	dq           queue.DelayQueue[TimingWheelSlot] // alignment 8, size 16; Do not use the timer.Ticker
	tasksMap     map[JobID]Task                    // alignment 8, size 8
	stopC        chan struct{}                     // alignment 8, size 8
	addTaskC     chan Task                         // alignment 8, size 8
	expiredSlotC chan TimingWheelSlot              // alignment 8, size 8
	taskCounter  *atomic.Int64                     // alignment 8, size 8
	slotCounter  *atomic.Int64                     // alignment 8, size 8
	cancelTaskC  chan JobID                        // alignment 8, size 8
	lock         *sync.RWMutex                     // alignment 8, size 8
	isRunning    *atomic.Bool                      // alignment 8, size 8
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
		addTaskC:     make(chan Task),
		stopC:        make(chan struct{}),
		expiredSlotC: make(chan TimingWheelSlot, 3),
		cancelTaskC:  make(chan JobID),
		tasksMap:     make(map[JobID]Task),
		lock:         &sync.RWMutex{},
		isRunning:    &atomic.Bool{},
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
	// It seems so stupid, but it is data race free.
	oldMap := xtw.tasksMap
	xtw.tasksMap = make(map[JobID]Task)
	clear(oldMap)

	close(xtw.stopC)
	close(xtw.addTaskC)
	close(xtw.expiredSlotC)
	close(xtw.cancelTaskC)
}

func (xtw *xTimingWheels) AddTask(task Task) error {
	if xtw.isRunning.Load() == false {
		return ErrTimingWheelStopped
	}
	if len(task.GetJobID()) <= 0 {
		return ErrTimingWheelTaskEmptyJobID
	}
	if task.GetJob() == nil {
		return ErrTimingWheelEmptyJob
	}
	if xtw.isRunning.Load() == false {
		return ErrTimingWheelStopped
	}
	xtw.addTaskC <- task // FIXME Block the caller
	return nil
}

func (xtw *xTimingWheels) AfterFunc(delayMs time.Duration, fn Job) (Task, error) {
	if xtw.isRunning.Load() == false {
		return nil, ErrTimingWheelStopped
	}
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
	if err := xtw.AddTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (xtw *xTimingWheels) ScheduleFunc(sched Scheduler, fn Job) (Task, error) {
	if xtw.isRunning.Load() == false {
		return nil, ErrTimingWheelStopped
	}
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
	if err := xtw.AddTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (xtw *xTimingWheels) CancelTask(jobID JobID) error {
	if xtw.isRunning.Load() == false {
		return ErrTimingWheelStopped
	}
	if len(jobID) <= 0 {
		return ErrTimingWheelTaskEmptyJobID
	}
	xtw.lock.RLock()
	task, ok := xtw.tasksMap[jobID]
	xtw.lock.RUnlock()
	if !ok {
		return ErrTimingWheelTaskNotFound
	}

	xtw.cancelTaskC <- task.GetJobID()
	return nil
}

func (xtw *xTimingWheels) schedule(ctx context.Context) {
	if ctx == nil {
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("[timing wheel] add task schedule panic recover", "error", err, "stack", debug.Stack())
			}
		}()
		for {
			select {
			case <-ctx.Done():
				xtw.Shutdown()
				return
			case <-xtw.stopC:
				return
			case task, ok := <-xtw.addTaskC:
				if !ok {
					slog.Warn("[timing wheel] add task channel is closed")
					continue
				}
				err := xtw.addTask(task)
				if errors.Is(err, ErrTimingWheelTaskIsExpired) {
					xtw.addOrRunTask(task)
				}
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
				xtw.lock.Lock()
				xtw.advanceClock(slot.GetExpirationMs())
				xtw.lock.Unlock()
				// Here related to slot level upgrade and downgrade.
				slot.Flush(xtw.addOrRunTask)
			}
		}
	}(ctx.Value(disableTimingWheelsScheduleExpiredSlot))
	go func(disabled any) {
		if disabled != nil && disabled.(bool) {
			return
		}
		defer func() {
			if err := recover(); err != nil {
				slog.Error("[timing wheel] cancel task schedule panic recover", "error", err, "stack", debug.Stack())
			}
		}()
		for {
			select {
			case <-ctx.Done():
				xtw.Shutdown()
				return
			case <-xtw.stopC:
				return
			case jobID, ok := <-xtw.cancelTaskC:
				if !ok {
					slog.Warn("[timing wheel] cancel task channel is closed")
					continue
				}
				xtw.cancelTask(jobID)
			}
		}
	}(ctx.Value(disableTimingWheelsScheduleCancelTask))
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
	if task == nil || task.Cancelled() || xtw.isRunning.Load() == false {
		return ErrTimingWheelStopped
	}
	// FIXME Recursive function to add a task, need to measure the performance.
	err := xtw.tw.(*timingWheel).addTask(task, 0)
	if err == nil || errors.Is(err, ErrTimingWheelTaskIsExpired) {
		xtw.lock.Lock()
		xtw.tasksMap[task.GetJobID()] = task
		xtw.lock.Unlock()
	}
	return err
}

func (xtw *xTimingWheels) addOrRunTask(t Task) {
	if t == nil || t.Cancelled() || xtw.isRunning.Load() == false {
		return
	}

	// FIXME goroutine pool to run this.
	// [slotExpMs, slotExpMs+interval)
	var (
		taskLevel int64
		runNow    bool
	)
	if t.GetSlot() != nil {
		taskLevel = t.GetSlot().GetLevel()
		runNow = t.GetSlot().GetExpirationMs() == sentinelSlotExpiredMs || taskLevel == 0
	} else {
		runNow = t.GetExpiredMs() <= time.Now().UTC().UnixMilli()
	}
	if runNow {
		go t.GetJob()(xtw.ctx, t.GetJobMetadata())
	}

	// Re-add loop job to timing wheel.
	// Upgrade and downgrade (move) the t from one slot to another slot.
	// Lock free.
	switch t.GetJobType() {
	case OnceJob:
		if taskLevel == 0 && runNow {
			xtw.cancelTask(t.GetJobID())
		} else if taskLevel > 0 && !runNow && xtw.isRunning.Load() == true {
			xtw.taskCounter.Add(-1)
			xtw.addTaskC <- t
		}
	case RepeatJob:
		var sTask Task
		if taskLevel == 0 && runNow {
			if t.GetRestLoopCount() == 0 {
				xtw.cancelTask(t.GetJobID())
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
			xtw.addTaskC <- sTask
		}
	}
	return
}

func (xtw *xTimingWheels) cancelTask(jobID JobID) {
	if xtw.isRunning.Load() == false {
		return
	}

	xtw.lock.RLock()
	task, ok := xtw.tasksMap[jobID]
	xtw.lock.RUnlock()
	if !ok {
		return
	}

	if task.GetSlot() != nil && !task.GetSlot().RemoveTask(task) {
		return
	}
	task.Cancel()

	xtw.lock.Lock()
	delete(xtw.tasksMap, jobID)
	xtw.lock.Unlock()
	xtw.taskCounter.Add(-1)
}
