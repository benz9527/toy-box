package pubsub

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestXSinglePipelineDisruptor(t *testing.T) {
	gTotal := 1000
	tasks := 100
	counter := &atomic.Int64{}
	wg := &sync.WaitGroup{}
	wg.Add(gTotal)
	disruptor := NewXSinglePipelineDisruptor[int](4,
		NewXNoCacheChannelBlockStrategy(),
		func(event int) error {
			counter.Add(1)
			return nil
		},
	)
	if err := disruptor.Start(); err != nil {
		t.Fatalf("disruptor start failed, err: %v", err)
	}
	for i := 0; i < gTotal; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < tasks; j++ {
				if _, _, err := disruptor.Publish(j); err != nil {
					t.Logf("publish failed, err: %v", err)
					break
				}
			}
		}()
	}
	wg.Wait()
	time.Sleep(time.Second)
	assert.Equal(t, int64(gTotal*tasks), counter.Load())
	err := disruptor.Stop()
	assert.NoError(t, err)
}
