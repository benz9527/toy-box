package queue

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

func TestXRingBufferCursor(t *testing.T) {
	cursor := NewXRingBufferCursor()
	beginIs := time.Now()
	for i := 0; i < 100000000; i++ {
		x := cursor.Increase()
		if x%10000000 == 0 {
			t.Logf("x=%d", x)
		}
	}
	diff := time.Since(beginIs)
	t.Logf("ts diff=%v", diff)
}

func TestXRingBufferCursorConcurrency(t *testing.T) {
	// lower than single goroutine test
	cursor := NewXRingBufferCursor()
	t.Logf("cursor size=%v", unsafe.Sizeof(*cursor.(*xRingBufferCursor)))
	beginIs := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func(idx int) {
			for j := 0; j < 10000; j++ {
				x := cursor.Increase()
				if x%10000000 == 0 {
					t.Logf("gid=%d, x=%d", idx, x)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	diff := time.Since(beginIs)
	t.Logf("ts diff=%v", diff)
}

func TestXRingBufferCursorNoPaddingConcurrency(t *testing.T) {
	// Better than padding version
	var cursor uint64 // same address, meaningless for data race
	beginIs := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func(idx int) {
			for j := 0; j < 10000; j++ {
				x := atomic.AddUint64(&cursor, 1)
				if x%10000000 == 0 {
					t.Logf("gid=%d, x=%d", idx, x)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	diff := time.Since(beginIs)
	t.Logf("ts diff=%v", diff)
}

type noPadObj struct {
	a, b, c uint64
}

func (o *noPadObj) increase() {
	atomic.AddUint64(&o.a, 1)
	atomic.AddUint64(&o.b, 1)
	atomic.AddUint64(&o.c, 1)
}

type padObj struct {
	a uint64
	_ [8]uint64
	b uint64
	_ [8]uint64
	c uint64
	_ [8]uint64
}

func (o *padObj) increase() {
	atomic.AddUint64(&o.a, 1)
	atomic.AddUint64(&o.b, 1)
	atomic.AddUint64(&o.c, 1)
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/queue
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkNoPadObj
// BenchmarkNoPadObj-4   	15984717	        79.60 ns/op
func BenchmarkNoPadObj(b *testing.B) {
	// Lower than padding version
	obj := &noPadObj{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj.increase()
		}
	})
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/queue
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkPadObj
// BenchmarkPadObj-4   	19536519	        62.11 ns/op
func BenchmarkPadObj(b *testing.B) {
	obj := &padObj{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj.increase()
		}
	})
}

func TestFalseSharing(t *testing.T) {
	// In the same cache line
	volatileArray := [4]uint64{0, 0, 0, 0} // contiguous memory
	var wg sync.WaitGroup
	wg.Add(4)
	beginIs := time.Now()
	for i := 0; i < 4; i++ {
		// Concurrent write to the same cache line
		// Many cache misses, because of many RFO
		go func(idx int) {
			for j := 0; j < 100000000; j++ {
				atomic.AddUint64(&volatileArray[idx], 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	diff := time.Since(beginIs)
	// ts diff=8.423525377s
	t.Logf("ts diff=%v", diff)
}

func TestNoFalseSharing(t *testing.T) {
	type padE struct {
		value uint64
		_     [7]uint64
	}

	// Each one in a different cache line
	volatileArray := [4]padE{{}, {}, {}, {}}
	var wg sync.WaitGroup
	wg.Add(4)
	beginIs := time.Now()
	for i := 0; i < 4; i++ {
		// No RFO data race
		go func(idx int) {
			for j := 0; j < 100000000; j++ {
				atomic.AddUint64(&volatileArray[idx].value, 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	diff := time.Since(beginIs)
	// ts diff=890.219393ms
	t.Logf("ts diff=%v", diff)
}
