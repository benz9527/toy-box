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

type timingWheel struct {
	ctx context.Context

	tickMs               int64
	startMs              int64 // baseline startup timestamp
	interval             int64
	currentTimeMs        int64
	slotSize             int64 // in kafka it is wheelSize
	globalTaskCounterRef *atomic.Int64
	globalSlotCounterRef *atomic.Int64

	lock *sync.RWMutex

	overflowWheelRef TimingWheel // same as kafka

	slots       []TimingWheelSlot // in kafka it is buckets
	globalDqRef queue.DelayQueue[TimingWheelSlot]
}

var (
	_ TimingWheel = (*timingWheel)(nil)
)

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
		currentTimeMs:        startMs - (startMs % tickMs),
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

func (tw *timingWheel) advanceClock(slotExpiredMs int64) {
	// Here related to slot level upgrade and downgrade.
	currentTimeMs := tw.GetCurrentTimeMs()
	tickMs := tw.GetTickMs()
	if slotExpiredMs >= currentTimeMs+tickMs {
		currentTimeMs = slotExpiredMs - (slotExpiredMs % tickMs)
		atomic.StoreInt64(&tw.currentTimeMs, currentTimeMs)
		if tw.overflowWheelRef != nil {
			tw.overflowWheelRef.(*timingWheel).advanceClock(currentTimeMs)
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

	taskExpiredMs := task.GetExpirationMs()
	currentTimeMs := tw.GetCurrentTimeMs()
	tickMs := tw.GetTickMs()
	interval := tw.GetInterval()
	slotSize := tw.GetSlotSize()

	if task.Cancelled() {
		return errors.New("[timing wheel] task is cancelled")
	} else if taskExpiredMs < currentTimeMs+tickMs {
		task.setSlot(immediateExpiredSlot)
		tw.globalSlotCounterRef.Add(1)
		return fmt.Errorf("[timing wheel] task taskExpiredMs at %d is before %d %w",
			taskExpiredMs, currentTimeMs+tickMs, ErrTimingWheelTaskIsExpired)
	} else if taskExpiredMs < currentTimeMs+interval {
		virtualID := taskExpiredMs / tickMs
		slotID := int(virtualID % slotSize)
		slot := tw.slots[slotID]
		if !slot.setExpirationMs(virtualID * tickMs) {
			return fmt.Errorf("[timing wheel] slot %d unable update the expiration %w",
				virtualID*tickMs, ErrTimingWheelTaskUnableToBeAddedToSlot)
		}

		tw.lock.Lock()
		slot.AddTask(task)
		slot.setSlotID(int64(slotID))
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

type xTimingWheels struct {
	tw TimingWheel

	ctx          context.Context
	stopC        chan struct{}
	addTaskC     chan Task
	expiredSlotC chan TimingWheelSlot
	taskCounter  *atomic.Int64
	slotCounter  *atomic.Int64
	cancelTaskC  chan JobID
	lock         *sync.RWMutex
	isRunning    atomic.Bool

	dq       queue.DelayQueue[TimingWheelSlot] // Do not use the timer.Ticker
	tasksMap map[JobID]Task
}

var (
	_ TimingWheels = (*xTimingWheels)(nil)
)

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
	xtw.dq = queue.NewArrayDelayQueue[TimingWheelSlot](32)
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

func (xtw *xTimingWheels) Shutdown() {
	if old := xtw.isRunning.Swap(false); !old {
		slog.Warn("[timing wheel] timing wheel is already shutdown")
		return
	}
	xtw.dq = nil
	xtw.isRunning.Store(false)
	oldMap := xtw.tasksMap
	clear(oldMap)
	xtw.tasksMap = make(map[JobID]Task)
	close(xtw.stopC)
	close(xtw.addTaskC)
	close(xtw.expiredSlotC)
	close(xtw.cancelTaskC)
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
	xtw.addTaskC <- task
	return nil
}

func (xtw *xTimingWheels) AfterFunc(delayMs time.Duration, fn func()) (Task, error) {
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
		func(ctx context.Context, jobID JobID) {
			fn()
		})

	if err := xtw.AddTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (xtw *xTimingWheels) ScheduleFunc(sched Scheduler, fn func()) (Task, error) {
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
		func(ctx context.Context, jobID JobID) {
			fn()
		})
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

const (
	disableTimingWheelsSchedulePoll        = "disableTWSPoll"
	disableTimingWheelsScheduleCancelTask  = "disableTWSCancelTask"
	disableTimingWheelsScheduleExpiredSlot = "disableTWSExpSlot"
)

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
					slog.Info("[timing wheel] task is expired immediately", "jobID", task.GetJobID())
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

func (xtw *xTimingWheels) advanceClock(timeoutMs int64) {
	// Here related to slot level upgrade and downgrade.
	xtw.tw.(*timingWheel).advanceClock(timeoutMs)
}

func (xtw *xTimingWheels) addTask(task Task) error {
	// FIXME Recursive function to add a task, need to measure the performance.
	err := xtw.tw.(*timingWheel).addTask(task, 0)
	if err == nil || errors.Is(err, ErrTimingWheelTaskIsExpired) {
		xtw.lock.Lock()
		xtw.tasksMap[task.GetJobID()] = task
		xtw.lock.Unlock()
	}
	return err
}

func (xtw *xTimingWheels) addOrRunTask(task Task) {
	if task == nil || task.Cancelled() {
		return
	}

	// FIXME goroutine pool to run this.
	// [slotExpMs, slotExpMs+interval)
	taskLevel := task.GetSlot().GetLevel()
	runNow := task.GetSlot().GetExpirationMs() == sentinelSlotExpiredMs || taskLevel == 0
	if runNow {
		go task.GetJob()(xtw.ctx, task.GetJobID())
	}

	// Re-add loop job to timing wheel.
	// Upgrade and downgrade (move) the task from one slot to another slot.
	switch task.GetTaskType() {
	case OnceTask:
		// FIXME lock race
		if taskLevel == 0 && runNow {
			xtw.cancelTask(task.GetJobID())
		} else if taskLevel > 0 && !runNow {
			xtw.taskCounter.Add(-1)
			// Lock free
			xtw.addTaskC <- task
		}
	case RepeatTask:
		var sTask Task
		if taskLevel == 0 && runNow {
			if task.GetRestLoopCount() == 0 {
				xtw.cancelTask(task.GetJobID())
				return
			}
			sTask, ok := task.(ScheduledTask)
			if !ok {
				return
			}
			sTask.UpdateNextScheduledMs()
			if sTask.GetExpirationMs() < 0 {
				return
			}
		} else if taskLevel > 0 && !runNow {
			sTask = task
		}
		// Lock free
		if sTask != nil {
			xtw.taskCounter.Add(-1)
			xtw.addTaskC <- sTask
		}
	}
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
