package timer

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"sort"
	"testing"
	"time"
)

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
		var (
			task Task
			err  error
		)
		task, err = tw.AfterFunc(delays[i], func(idx int) func() {
			return func() {
				execAt := time.Now().UTC().UnixMilli()
				slog.Info("after func", "delay", delays[idx].String(), "exec at", execAt, "diff",
					execAt-tw.GetStartMs()-delays[idx].Milliseconds())
			}
		}(i))
		assert.NoError(t, err)
		t.Logf("task: %s\n", task.GetJobID())
	}
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
		_, err := tw.AfterFunc(delays[i], func(idx int) func() {
			return func() {}
		}(i))
		if err != nil {

		}
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
		t.Logf("slot level: %d, ID %d, %d\n", v.GetSlot().GetLevel(), v.GetSlot().GetSlotID(), v.GetExpirationMs())
	}
	<-ctx.Done()
}
