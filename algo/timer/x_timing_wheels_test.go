package timer

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewTimingWheels(t *testing.T) {
	ctx, cancel := context.WithTimeoutCause(context.TODO(), 2*time.Second, errors.New("timeout"))
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
	t.Logf("tw taskCounter: %d\n", tw.GetTaskCounter())
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)
}
