package queue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"unsafe"
)

func TestDelayQueueAlignmentAndSize(t *testing.T) {
	dq := NewArrayDelayQueue[*person](32)
	t.Logf("dq aligment size: %d\n", unsafe.Alignof(dq))
	prototype := dq.(*arrayDQ[*person])
	t.Logf("dq prototype alignment size: %d\n", unsafe.Alignof(prototype))
	t.Logf("dq prototype item channel alignment size: %d\n", unsafe.Alignof(prototype.itemC))
	t.Logf("dq prototype wake up channel alignment size: %d\n", unsafe.Alignof(prototype.wakeUpC))
	t.Logf("dq prototype priority queue alignment size: %d\n", unsafe.Alignof(prototype.pq))
	t.Logf("dq prototype sleeping alignment size: %d\n", unsafe.Alignof(prototype.sleeping))
	t.Logf("dq prototype lock alignment size: %d\n", unsafe.Alignof(prototype.lock))

	t.Logf("dq size: %d\n", unsafe.Sizeof(dq))
	t.Logf("dq prototype size: %d\n", unsafe.Sizeof(prototype))
	t.Logf("dq prototype item channel size: %d\n", unsafe.Sizeof(prototype.itemC))
	t.Logf("dq prototype wake up channel size: %d\n", unsafe.Sizeof(prototype.wakeUpC))
	t.Logf("dq prototype priority queue size: %d\n", unsafe.Sizeof(prototype.pq))
	t.Logf("dq prototype sleeping size: %d\n", unsafe.Sizeof(prototype.sleeping))
	t.Logf("dq prototype lock size: %d\n", unsafe.Sizeof(prototype.lock))

	assert.Equal(t, uintptr(48), unsafe.Sizeof(*prototype))
}

func TestDelayQueue_Poll(t *testing.T) {
	ms := time.Now().UnixMilli()
	dq := NewArrayDelayQueue[*person](32)
	_ = dq.Offer(&person{age: 10, name: "p0", salary: ms + 110}, ms+110)
	_ = dq.Offer(&person{age: 101, name: "p1", salary: ms + 501}, ms+501)
	_ = dq.Offer(&person{age: 10, name: "p2", salary: ms + 155}, ms+155)
	_ = dq.Offer(&person{age: 200, name: "p3", salary: ms + 210}, ms+210)
	_ = dq.Offer(&person{age: 3, name: "p4", salary: ms + 60}, ms+60)
	_ = dq.Offer(&person{age: 1, name: "p5", salary: ms + 110}, ms+110)
	_ = dq.Offer(&person{age: 5, name: "p6", salary: ms + 250}, ms+250)
	_ = dq.Offer(&person{age: 200, name: "p7", salary: ms + 301}, ms+301)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go dq.Poll(ctx,
		func() int64 {
			return time.Now().UnixMilli()
		},
	)
	for {
		select {
		case item := <-dq.Wait():
			t.Logf("current time ms: %d, item: %v\n", time.Now().UnixMilli(), item)
		case <-ctx.Done():
			return
		}
	}
}