package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTimingWheelEventsPool(t *testing.T) {
	pool := newTimingWheelEventsPool()
	if pool == nil {
		t.Fatal("pool is nil")
	}
	event := pool.Get()
	event.AddTask(&task{
		jobMetadata: &jobMetadata{
			jobID: "1",
		},
	})
	task, ok := event.GetAddTask()
	assert.True(t, ok)
	assert.Equal(t, JobID("1"), task.GetJobID())
	jobID, ok := event.GetCancelTaskJobID()
	assert.False(t, ok)
	pool.Put(event)

	event = pool.Get()
	event.CancelTaskJobID("2")
	jobID, ok = event.GetCancelTaskJobID()
	assert.True(t, ok)
	assert.Equal(t, JobID("2"), jobID)
}
