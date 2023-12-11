package timer

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"sort"
	"testing"
	"time"
	"unsafe"
)

func TestTimingWheel_AlignmentAndSize(t *testing.T) {
	tw := &timingWheel{}
	t.Logf("tw alignment: %d\n", unsafe.Alignof(tw))
	t.Logf("tw ctx alignment: %d\n", unsafe.Alignof(tw.ctx))
	t.Logf("tw slot alignment: %d\n", unsafe.Alignof(tw.slots))
	t.Logf("tw tickMs alignment: %d\n", unsafe.Alignof(tw.tickMs))
	t.Logf("tw startMs alignment: %d\n", unsafe.Alignof(tw.startMs))
	t.Logf("tw slotSize alignment: %d\n", unsafe.Alignof(tw.slotSize))
	t.Logf("tw globalTaskCounterRef alignment: %d\n", unsafe.Alignof(tw.globalTaskCounterRef))
	t.Logf("tw globalSlotCounterRef alignment: %d\n", unsafe.Alignof(tw.globalSlotCounterRef))
	t.Logf("tw overflowWheelRef alignment: %d\n", unsafe.Alignof(tw.overflowWheelRef))
	t.Logf("tw globalDqRef alignment: %d\n", unsafe.Alignof(tw.globalDqRef))

	t.Logf("tw size: %d\n", unsafe.Sizeof(*tw))
	t.Logf("tw ctx size: %d\n", unsafe.Sizeof(tw.ctx))
	t.Logf("tw slot size: %d\n", unsafe.Sizeof(tw.slots))
	t.Logf("tw tickMs size: %d\n", unsafe.Sizeof(tw.tickMs))
	t.Logf("tw startMs size: %d\n", unsafe.Sizeof(tw.startMs))
	t.Logf("tw slotSize size: %d\n", unsafe.Sizeof(tw.slotSize))
	t.Logf("tw globalTaskCounterRef size: %d\n", unsafe.Sizeof(tw.globalTaskCounterRef))
	t.Logf("tw globalSlotCounterRef size: %d\n", unsafe.Sizeof(tw.globalSlotCounterRef))
	t.Logf("tw overflowWheelRef size: %d\n", unsafe.Sizeof(tw.overflowWheelRef))
	t.Logf("tw globalDqRef size: %d\n", unsafe.Sizeof(tw.globalDqRef))
}

func TestXTimingWheels_AlignmentAndSize(t *testing.T) {
	tw := &xTimingWheels{}
	t.Logf("tw alignment: %d\n", unsafe.Alignof(tw))
	t.Logf("tw ctx alignment: %d\n", unsafe.Alignof(tw.ctx))
	t.Logf("tw tw alignment: %d\n", unsafe.Alignof(tw.tw))
	t.Logf("tw stopC alignment: %d\n", unsafe.Alignof(tw.stopC))
	t.Logf("tw twEventC alignment: %d\n", unsafe.Alignof(tw.twEventC))
	t.Logf("tw twEventPool alignment: %d\n", unsafe.Alignof(tw.twEventPool))
	t.Logf("tw expiredSlotC alignment: %d\n", unsafe.Alignof(tw.expiredSlotC))
	t.Logf("tw taskCounter alignment: %d\n", unsafe.Alignof(tw.taskCounter))
	t.Logf("tw slotCounter alignment: %d\n", unsafe.Alignof(tw.slotCounter))
	t.Logf("tw isRunning alignment: %d\n", unsafe.Alignof(tw.isRunning))
	t.Logf("tw dq alignment: %d\n", unsafe.Alignof(tw.dq))
	t.Logf("tw tasksMap alignment: %d\n", unsafe.Alignof(tw.tasksMap))

	t.Logf("tw size: %d\n", unsafe.Sizeof(*tw))
	t.Logf("tw ctx size: %d\n", unsafe.Sizeof(tw.ctx))
	t.Logf("tw tw size: %d\n", unsafe.Sizeof(tw.tw))
	t.Logf("tw stopC size: %d\n", unsafe.Sizeof(tw.stopC))
	t.Logf("tw twEventC size: %d\n", unsafe.Sizeof(tw.twEventC))
	t.Logf("tw twEventPool size: %d\n", unsafe.Sizeof(tw.twEventPool))
	t.Logf("tw expiredSlotC size: %d\n", unsafe.Sizeof(tw.expiredSlotC))
	t.Logf("tw taskCounter size: %d\n", unsafe.Sizeof(tw.taskCounter))
	t.Logf("tw slotCounter size: %d\n", unsafe.Sizeof(tw.slotCounter))
	t.Logf("tw isRunning size: %d\n", unsafe.Sizeof(tw.isRunning))
	t.Logf("tw dq size: %d\n", unsafe.Sizeof(tw.dq))
	t.Logf("tw tasksMap size: %d\n", unsafe.Sizeof(tw.tasksMap))
}

func TestNewTimingWheels(t *testing.T) {
	ctx, cancel := context.WithTimeoutCause(context.TODO(), time.Second, errors.New("timeout"))
	defer cancel()
	tw := NewTimingWheels(
		ctx,
		time.Now().UTC().UnixMilli(),
		WithTimingWheelTickMs(100),
		WithTimingWheelSlotSize(32),
	)
	t.Logf("tw tickMs: %d\n", tw.GetTickMs())
	t.Logf("tw startMs: %d\n", tw.GetStartMs())
	t.Logf("tw slotSize: %d\n", tw.GetSlotSize())
	t.Logf("tw globalTaskCounterRef: %d\n", tw.GetTaskCounter())
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)
}

func TestXTimingWheels_AfterFunc(t *testing.T) {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 2*time.Second, errors.New("timeout"))
	defer cancel()
	tw := NewTimingWheels(
		ctx,
		time.Now().UTC().UnixMilli(),
	)

	delays := []time.Duration{
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		15 * time.Millisecond,
		18 * time.Millisecond,
		20 * time.Millisecond,
		21 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		200 * time.Millisecond,
		500 * time.Millisecond,
		time.Second,
	}

	for i := 0; i < len(delays); i++ {
		_, err := tw.AfterFunc(delays[i], func(ctx context.Context, md JobMetadata) {
			execAt := time.Now().UTC().UnixMilli()
			slog.Info("after func", "expired ms", md.GetExpiredMs(), "exec at", execAt, "diff",
				execAt-md.GetExpiredMs())
		})
		assert.NoError(t, err)
	}
	t.Logf("tw tickMs: %d\n", tw.GetTickMs())
	t.Logf("tw startMs: %d\n", tw.GetStartMs())
	t.Logf("tw slotSize: %d\n", tw.GetSlotSize())
	t.Logf("tw tasks: %d\n", tw.GetTaskCounter())
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)
	t.Logf("final tw tasks: %d\n", tw.GetTaskCounter())
}

func TestXTimingWheels_ScheduleFunc_ConcurrentFinite(t *testing.T) {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 3*time.Second, errors.New("timeout"))
	defer cancel()
	tw := NewTimingWheels(
		ctx,
		time.Now().UTC().UnixMilli(),
	)

	delays := []time.Duration{
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		15 * time.Millisecond,
		18 * time.Millisecond,
		20 * time.Millisecond,
		21 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		200 * time.Millisecond,
		500 * time.Millisecond,
		time.Second,
	}
	sched := NewFiniteScheduler(delays...)
	assert.NotNil(t, sched)
	// FIXME Here will occur some errors:
	//  1. only one task will be scheduled
	//  2. one scheduled task will lost in the timing wheels
	//  3. sometime the scheduling will be blocked
	go func() {
		task, err := tw.ScheduleFunc(sched, func(ctx context.Context, md JobMetadata) {
			execAt := time.Now().UTC().UnixMilli()
			slog.Info("sched1 after func", "expired ms", md.GetExpiredMs(), "exec at", execAt, "diff",
				execAt-md.GetExpiredMs())
		})
		assert.NoError(t, err)
		t.Logf("task1: %s\n", task.GetJobID())
	}()
	go func() {
		task, err := tw.ScheduleFunc(sched, func(ctx context.Context, md JobMetadata) {
			execAt := time.Now().UTC().UnixMilli()
			slog.Info("sched2 after func", "expired ms", md.GetExpiredMs(), "exec at", execAt, "diff",
				execAt-md.GetExpiredMs())
		})
		assert.NoError(t, err)
		t.Logf("task2: %s\n", task.GetJobID())
	}()

	t.Logf("tw tickMs: %d\n", tw.GetTickMs())
	t.Logf("tw startMs: %d\n", tw.GetStartMs())
	t.Logf("tw slotSize: %d\n", tw.GetSlotSize())
	t.Logf("tw tasks: %d\n", tw.GetTaskCounter())
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)
	t.Logf("final tw tasks: %d\n", tw.GetTaskCounter())
}

func TestXTimingWheels_ScheduleFunc_1MsInfinite(t *testing.T) {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 1*time.Second, errors.New("timeout"))
	defer cancel()
	tw := NewTimingWheels(
		ctx,
		time.Now().UTC().UnixMilli(),
	)

	delays := []time.Duration{
		time.Millisecond,
	}
	sched := NewInfiniteScheduler(delays...)
	assert.NotNil(t, sched)
	// FIXME 1ms infinite scheduling will occur critical delay execution error
	task, err := tw.ScheduleFunc(sched, func(ctx context.Context, md JobMetadata) {
		execAt := time.Now().UTC().UnixMilli()
		slog.Info("infinite sched1 after func", "expired ms", md.GetExpiredMs(), "exec at", execAt, "diff",
			execAt-md.GetExpiredMs())
	})
	assert.NoError(t, err)
	t.Logf("task1: %s\n", task.GetJobID())

	t.Logf("tw tickMs: %d\n", tw.GetTickMs())
	t.Logf("tw startMs: %d\n", tw.GetStartMs())
	t.Logf("tw slotSize: %d\n", tw.GetSlotSize())
	t.Logf("tw tasks: %d\n", tw.GetTaskCounter())
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)
	t.Logf("final tw tasks: %d\n", tw.GetTaskCounter())
}

func TestXTimingWheels_ScheduleFunc_32MsInfinite(t *testing.T) {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 3*time.Second, errors.New("timeout"))
	defer cancel()
	tw := NewTimingWheels(
		ctx,
		time.Now().UTC().UnixMilli(),
	)

	delays := []time.Duration{
		18 * time.Millisecond,
	}
	sched := NewInfiniteScheduler(delays...)
	assert.NotNil(t, sched)
	task, err := tw.ScheduleFunc(sched, func(ctx context.Context, md JobMetadata) {
		execAt := time.Now().UTC().UnixMilli()
		slog.Info("infinite sched32 after func", "expired ms", md.GetExpiredMs(), "exec at", execAt, "diff",
			execAt-md.GetExpiredMs())
	})
	assert.NoError(t, err)
	t.Logf("task1: %s\n", task.GetJobID())

	t.Logf("tw tickMs: %d\n", tw.GetTickMs())
	t.Logf("tw startMs: %d\n", tw.GetStartMs())
	t.Logf("tw slotSize: %d\n", tw.GetSlotSize())
	t.Logf("tw tasks: %d\n", tw.GetTaskCounter())
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)
	t.Logf("final tw tasks: %d\n", tw.GetTaskCounter())
}

func TestXTimingWheels_AfterFunc_Slots(t *testing.T) {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 500*time.Millisecond, errors.New("timeout"))
	defer cancel()
	ctx = context.WithValue(ctx, disableTimingWheelsScheduleCancelTask, true)
	ctx = context.WithValue(ctx, disableTimingWheelsScheduleExpiredSlot, true)
	ctx = context.WithValue(ctx, disableTimingWheelsSchedulePoll, true)

	tw := NewTimingWheels(
		ctx,
		time.Now().UTC().UnixMilli(),
	)

	delays := []time.Duration{
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		15 * time.Millisecond,
		18 * time.Millisecond,
		20 * time.Millisecond,
		21 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
		time.Second,
	}

	for i := 0; i < len(delays); i++ {
		_, err := tw.AfterFunc(delays[i], func(ctx context.Context, md JobMetadata) {})
		assert.NoError(t, err)
	}
	t.Logf("tw tickMs: %d\n", tw.GetTickMs())
	t.Logf("tw startMs: %d\n", tw.GetStartMs())
	t.Logf("tw slotSize: %d\n", tw.GetSlotSize())
	t.Logf("tw tasks: %d\n", tw.GetTaskCounter())

	<-time.After(100 * time.Millisecond)
	taskIDs := make([]string, 0, len(delays))
	for k := range tw.(*xTimingWheels).tasksMap {
		taskIDs = append(taskIDs, string(k))
	}
	sort.Strings(taskIDs)
	for _, taskID := range taskIDs {
		v := tw.(*xTimingWheels).tasksMap[JobID(taskID)]
		t.Logf("slot level: %d, ID %d, %d\n", v.GetSlot().GetLevel(), v.GetSlot().GetSlotID(), v.GetExpiredMs())
	}
	<-ctx.Done()
}

func BenchmarkNewTimingWheels_AfterFunc(b *testing.B) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, disableTimingWheelsScheduleCancelTask, true)
	ctx = context.WithValue(ctx, disableTimingWheelsScheduleExpiredSlot, true)
	ctx = context.WithValue(ctx, disableTimingWheelsSchedulePoll, true)
	tw := NewTimingWheels(
		ctx,
		time.Now().UTC().UnixMilli(),
		WithTimingWheelTickMs(1),
		WithTimingWheelSlotSize(20),
	)
	defer tw.Shutdown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tw.AfterFunc(time.Duration(i+1)*time.Millisecond, func(ctx context.Context, md JobMetadata) {
		})
		assert.NoError(b, err)
	}
	b.ReportAllocs()
}
