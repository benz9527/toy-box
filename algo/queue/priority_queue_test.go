package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

type person struct {
	name   string
	age    int
	salary int64
}

func TestPriorityQueueItemAlignmentAndSize(t *testing.T) {
	item := NewPQItem[*person](&person{age: 10, name: "p0"}, 1)
	t.Logf("item alignment size: %d\n", unsafe.Alignof(item))
	prototype := item.(*pqItem[*person])
	t.Logf("item prototype alignment size: %d\n", unsafe.Alignof(prototype))
	t.Logf("item prototype value alignment size: %d\n", unsafe.Alignof(prototype.value))
	t.Logf("item prototype priority alignment size: %d\n", unsafe.Alignof(prototype.priority))
	t.Logf("item prototype index alignment size: %d\n", unsafe.Alignof(prototype.index))
	t.Logf("item prototype comparator alignment size: %d\n", unsafe.Alignof(prototype.comparator))
	t.Logf("item prototype size: %d\n", unsafe.Sizeof(prototype))
	t.Logf("item prototype value size: %d\n", unsafe.Sizeof(prototype.value))
	t.Logf("item prototype priority size: %d\n", unsafe.Sizeof(prototype.priority))
	t.Logf("item prototype index size: %d\n", unsafe.Sizeof(prototype.index))
	t.Logf("item prototype comparator size: %d\n", unsafe.Sizeof(prototype.comparator))
}

func TestPriorityQueue_MinValueAsHighPriority(t *testing.T) {
	pq := NewArrayPriorityQueue[*person](32,
		func(i, j PQItem[*person]) bool {
			return i.GetPriority() < j.GetPriority()
		},
	)
	pq.Push(NewPQItem[*person](&person{age: 10, name: "p0"}, 1))
	pq.Push(NewPQItem[*person](&person{age: 101, name: "p1"}, 101))
	pq.Push(NewPQItem[*person](&person{age: 10, name: "p2"}, 10))
	pq.Push(NewPQItem[*person](&person{age: 200, name: "p3"}, 200))
	pq.Push(NewPQItem[*person](&person{age: 3, name: "p4"}, 3))
	pq.Push(NewPQItem[*person](&person{age: 1, name: "p5"}, 1))
	pq.Push(NewPQItem[*person](&person{age: 5, name: "p6"}, 5))

	expectedPriorities := []int64{1, 1, 3, 5, 10, 101, 200}
	for i, priority := range expectedPriorities {
		item := pq.Pop()
		t.Logf("%v， priority: %d", item.GetValue(), item.GetPriority())
		assert.Equal(t, priority, item.GetPriority(), "priority", i)
	}
}

func TestPriorityQueue_MaxValueAsHighPriority(t *testing.T) {
	pq := NewArrayPriorityQueue[*person](32,
		func(i, j PQItem[*person]) bool {
			return i.GetPriority() > j.GetPriority()
		},
	)
	pq.Push(NewPQItem[*person](&person{age: 10, name: "p0"}, 1))
	pq.Push(NewPQItem[*person](&person{age: 101, name: "p1"}, 101))
	pq.Push(NewPQItem[*person](&person{age: 10, name: "p2"}, 10))
	pq.Push(NewPQItem[*person](&person{age: 200, name: "p3"}, 200))
	pq.Push(NewPQItem[*person](&person{age: 3, name: "p4"}, 3))
	pq.Push(NewPQItem[*person](&person{age: 1, name: "p5"}, 1))
	pq.Push(NewPQItem[*person](&person{age: 5, name: "p6"}, 5))
	pq.Push(NewPQItem[*person](&person{age: 200, name: "p7"}, 201))

	expectedPriorities := []int64{201, 200, 101, 10, 5, 3, 1, 1}
	for i, priority := range expectedPriorities {
		item := pq.Pop()
		t.Logf("%v， priority: %d", item.GetValue(), item.GetPriority())
		assert.Equal(t, priority, item.GetPriority(), "priority", i)
	}
}

func TestPriorityQueue_MaxValueAsHighPriority_Peek(t *testing.T) {
	pq := NewArrayPriorityQueue[*person](32,
		func(i, j PQItem[*person]) bool {
			return i.GetPriority() > j.GetPriority()
		},
	)
	pq.Push(NewPQItem[*person](&person{age: 10, name: "p0"}, 1))
	pq.Push(NewPQItem[*person](&person{age: 101, name: "p1"}, 101))
	pq.Push(NewPQItem[*person](&person{age: 10, name: "p2"}, 10))
	pq.Push(NewPQItem[*person](&person{age: 200, name: "p3"}, 200))
	pq.Push(NewPQItem[*person](&person{age: 3, name: "p4"}, 3))
	pq.Push(NewPQItem[*person](&person{age: 1, name: "p5"}, 1))
	pq.Push(NewPQItem[*person](&person{age: 5, name: "p6"}, 5))
	pq.Push(NewPQItem[*person](&person{age: 200, name: "p7"}, 201))

	expectedPriorities := []int64{201, 200, 101, 10, 5, 3, 1, 1}
	for i, priority := range expectedPriorities {
		peekItem := pq.Peek()
		t.Logf("peek item: %v， priority: %d", peekItem.GetValue(), peekItem.GetPriority())
		item := pq.Pop()
		t.Logf("%v， priority: %d", item.GetValue(), item.GetPriority())
		assert.Equal(t, priority, item.GetPriority(), "priority", i)
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/queue
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkPriorityQueue_Push
// BenchmarkPriorityQueue_Push-4   	 4362896	       229.8 ns/op	      96 B/op	       0 allocs/op
func BenchmarkPriorityQueue_Push(b *testing.B) {
	var list = make([]PQItem[*person], 0, b.N)
	for i := 0; i < b.N; i++ {
		e := NewPQItem[*person](&person{age: i, name: fmt.Sprintf("p%d", i)}, int64(i))
		list = append(list, e)
	}
	b.ResetTimer()
	pq := NewArrayPriorityQueue[*person](32,
		func(i, j PQItem[*person]) bool {
			return i.GetPriority() < j.GetPriority()
		},
	)
	for i := 0; i < b.N; i++ {
		pq.Push(list[i])
	}
	b.ReportAllocs()
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/queue
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkPriorityQueue_Pop
// BenchmarkPriorityQueue_Pop-4   	  818748	      2113 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPriorityQueue_Pop(b *testing.B) {
	var list = make([]PQItem[*person], 0, b.N)
	for i := 0; i < b.N; i++ {
		e := NewPQItem[*person](&person{age: i, name: fmt.Sprintf("p%d", i)}, int64(i))
		list = append(list, e)
	}
	pq := NewArrayPriorityQueue[*person](32)
	for i := 0; i < b.N; i++ {
		pq.Push(list[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pq.Pop()
	}
	b.ReportAllocs()
}
