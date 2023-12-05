package timer

import (
	"context"
	"errors"
	"fmt"
	"github.com/benz9527/toy-box/algo/queue"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type timingWheel struct {
	ctx context.Context

	tickMs        int64
	startMs       int64 // baseline startup timestamp
	interval      int64
	currentTimeMs int64
	slotSize      int64 // in kafka it is wheelSize
	taskCounter   *atomic.Int64

	lock *sync.RWMutex

	overflowWheelRef TimingWheel // same as kafka

	slots       []TimingWheelSlot // in kafka it is buckets
	globalDqRef queue.DelayQueue[TimingWheelSlot]
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
	dq queue.DelayQueue[TimingWheelSlot],
) TimingWheel {
	tw := &timingWheel{
		ctx:           ctx,
		lock:          &sync.RWMutex{},
		tickMs:        tickMs,
		startMs:       startMs,
		slotSize:      slotSize,
		taskCounter:   taskCounter,
		interval:      tickMs * slotSize,
		currentTimeMs: truncateMs(startMs, tickMs),
		slots:         make([]TimingWheelSlot, slotSize),
		globalDqRef:   dq,
	}
	// Slot initialize by doubly linked list.
	for i := int64(0); i < slotSize; i++ {
		tw.slots[i] = NewXSlot()
	}
	return tw
}

func (tw *timingWheel) GetTickMs() int64 {
	return atomic.LoadInt64(&tw.tickMs)
}

func (tw *timingWheel) GetStartMs() int64 {
	return atomic.LoadInt64(&tw.startMs)
}

func (tw *timingWheel) GetTaskCounter() int64 {
	return tw.taskCounter.Load()
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

func (tw *timingWheel) advanceClock(timeoutMs int64) {
	currentTimeMs := tw.GetCurrentTimeMs()
	tickMs := tw.GetTickMs()
	if timeoutMs >= currentTimeMs+tickMs {
		currentTimeMs = truncateMs(timeoutMs, tickMs)
		if tw.overflowWheelRef != nil {
			tw.overflowWheelRef.(*timingWheel).advanceClock(currentTimeMs)
		}
	}
}

func (tw *timingWheel) addTask(task Task) error {
	if len(task.GetJobID()) <= 0 {
		return ErrTimingWheelTaskEmptyJobID
	}
	if task.GetJob() == nil {
		return ErrTimingWheelEmptyJob
	}

	tw.lock.Lock()
	defer tw.lock.Unlock()

	taskExpiration := task.GetExpirationMs()
	currentTimeMs := tw.GetCurrentTimeMs()
	tickMs := tw.GetTickMs()
	interval := tw.GetInterval()
	slotSize := tw.GetSlotSize()

	switch {
	case task.Cancelled():
		return errors.New("[timing wheel] task is cancelled")
	case taskExpiration < currentTimeMs+tickMs:
		return fmt.Errorf("[timing wheel] task taskExpiration at %d is before %d %w",
			taskExpiration, currentTimeMs+tickMs, ErrTimingWheelTaskIsExpired)
	case taskExpiration < currentTimeMs+interval:
		virtualID := taskExpiration / tickMs
		slotID := int(virtualID % slotSize)
		slot := tw.slots[slotID]
		if !slot.SetExpirationMs(virtualID * tickMs) {
			return fmt.Errorf("[timing wheel] slot %d unable update the expiration %w",
				virtualID*tickMs, ErrTimingWheelTaskUnableToBeAddedToSlot)
		}
		slot.AddTask(task)
		_ = tw.globalDqRef.Offer(slot, slot.GetExpirationMs())
		tw.taskCounter.Add(1)
		return nil
	default:
		// Out of the interval. Put it into the higher interval timing wheel
		if tw.overflowWheelRef == nil {
			tw.overflowWheelRef = newTimingWheel(
				tw.ctx,
				interval,
				slotSize,
				currentTimeMs,
				tw.taskCounter,
				tw.globalDqRef,
			)
		}
		// Tail recursive call, it will be free the previous stack frame.
		return tw.overflowWheelRef.(*timingWheel).addTask(task)
	}
}

type xTimingWheels struct {
	tw TimingWheel

	ctx          context.Context
	stopC        chan struct{}
	addTaskC     chan Task
	expiredSlotC chan TimingWheelSlot
	taskCounter  *atomic.Int64
	cancelTaskC  chan JobID
	lock         sync.RWMutex
	isRunning    atomic.Bool

	dq       queue.DelayQueue[TimingWheelSlot] // Do not use the timer.Ticker
	tasksMap map[JobID]Task
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
		addTaskC:     make(chan Task),
		stopC:        make(chan struct{}),
		expiredSlotC: make(chan TimingWheelSlot),
		cancelTaskC:  make(chan JobID),
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

	if tw.taskCounter == nil {
		tw.taskCounter = &atomic.Int64{}
		tw.taskCounter.Store(0)
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
		xtw.dq,
	)
	xtw.startAndRun(ctx)
	return xtw
}

func (xtw *xTimingWheels) Shutdown() {
	if old := xtw.isRunning.Swap(false); !old {
		slog.Warn("[timing wheel] timing wheel is already shutdown")
		return
	}
	xtw.dq = nil
	xtw.tasksMap = nil
	xtw.isRunning.Store(false)
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
	return xtw.tw.GetSlotSize()
}

func (xtw *xTimingWheels) AddTask(task Task) error {
	if len(task.GetJobID()) <= 0 {
		return ErrTimingWheelTaskEmptyJobID
	}
	if task.GetJob() == nil {
		return ErrTimingWheelEmptyJob
	}
	xtw.addTaskC <- task
	return nil
}

func (xtw *xTimingWheels) CancelTask(jobID JobID) error {
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

func (xtw *xTimingWheels) startAndRun(ctx context.Context) {
	if ctx == nil {
		return
	}
	go func() {
		_ = xtw.dq.PollToChannel(xtw.ctx, func() int64 {
			return time.Now().UTC().UnixMilli()
		}, xtw.expiredSlotC)
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("[timing wheel] scheduler panic recover", "error", err)
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
					slog.Warn("[timing wheel] expired slot channel is closed")
					continue
				}
				xtw.advanceClock(slot.GetExpirationMs())
				slot.Flush(xtw.addOrRunTask)
			case task, ok := <-xtw.addTaskC:
				if !ok {
					slog.Warn("[timing wheel] add task channel is closed")
					continue
				}
				err := xtw.addTask(task)
				if errors.Is(err, ErrTimingWheelTaskIsExpired) {
					// Run task immediately
				}
			case jobID, ok := <-xtw.cancelTaskC:
				if !ok {
					slog.Warn("[timing wheel] cancel task channel is closed")
					continue
				}
				xtw.cancelTask(jobID)
			}
		}
	}()
	xtw.isRunning.Store(true)
}

func (xtw *xTimingWheels) advanceClock(timeoutMs int64) {
	xtw.tw.(*timingWheel).advanceClock(timeoutMs)
}

func (xtw *xTimingWheels) addTask(task Task) error {
	// FIXME Recursive function to add a task, need to measure the performance.
	err := xtw.tw.(*timingWheel).addTask(task)
	if err == nil {
		// FIXME thread safe?
		xtw.tasksMap[task.GetJobID()] = task
	}
	return err
}

func (xtw *xTimingWheels) addOrRunTask(task Task) {
	if task == nil || task.Cancelled() {
		return
	}
	// FIXME goroutine pool to run this.
	go task.GetJob()(context.TODO(), task.GetJobID())

	xtw.lock.Lock()
	defer xtw.lock.Unlock()

	// Re-add loop job to timing wheel.
	if task.GetLoopCount() >= 1 {
		// FIXME Re-add error handle
		_ = xtw.AddTask(task)
		task.DecreaseLoopCount()
	} else if task.GetLoopCount() <= -1 {
		_ = xtw.AddTask(task)
	}
}

func (xtw *xTimingWheels) cancelTask(jobID JobID) {
	task, ok := xtw.tasksMap[jobID]
	if !ok {
		return
	}

	xtw.lock.Lock()
	defer xtw.lock.Unlock()

	if !task.GetSlot().RemoveTask(task) {
		return
	}
	task.Cancel()
	delete(xtw.tasksMap, jobID)
	xtw.taskCounter.Add(-1)
}

func truncateMs(srcMs, baseMs int64) int64 {
	if baseMs <= 0 {
		return srcMs
	}
	return srcMs - (srcMs % baseMs)
}
