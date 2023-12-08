package list

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList_AppendValue(t *testing.T) {
	dlist := NewLinkedList[int]()
	elements := dlist.AppendValue(1, 2, 3, 4, 5)
	assert.Equal(t, len(elements), 5)
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, e: %v", idx, e)
		assert.Equal(t, elements[idx], e)
		t.Logf("addr: %p, return addr: %p", elements[idx], e)
	})

	dlist2 := list.New()
	dlist2.PushBack(1)
	dlist2.PushBack(2)
	dlist2.PushBack(3)
	dlist2.PushBack(4)
	dlist2.PushBack(5)

	assert.Equal(t, dlist.Len(), int64(dlist2.Len()))

	dlistItr := dlist.Front()
	dlist2Itr := dlist2.Front()
	for dlist2Itr != nil {
		assert.Equal(t, dlistItr.GetValue(), dlist2Itr.Value)
		dlist2Itr = dlist2Itr.Next()
		dlistItr = dlistItr.GetNext()
	}
}

func TestDoublyLinkedList_InsertBefore(t *testing.T) {
	dlist := NewLinkedList[int]()
	elements := dlist.AppendValue(1)
	_2n := dlist.InsertBefore(2, elements[0])
	_3n := dlist.InsertBefore(3, _2n)
	_4n := dlist.InsertBefore(4, _3n)
	dlist.InsertBefore(5, _4n)
	assert.Equal(t, int64(5), dlist.Len())

	dlist2 := list.New()
	_1n_2 := dlist2.PushBack(1)
	_2n_2 := dlist2.InsertBefore(2, _1n_2)
	_3n_2 := dlist2.InsertBefore(3, _2n_2)
	_4n_2 := dlist2.InsertBefore(4, _3n_2)
	dlist2.InsertBefore(5, _4n_2)

	assert.Equal(t, dlist.Len(), int64(dlist2.Len()))

	dlistItr := dlist.Front()
	dlist2Itr := dlist2.Front()
	for dlist2Itr != nil {
		assert.Equal(t, dlistItr.GetValue(), dlist2Itr.Value)
		dlist2Itr = dlist2Itr.Next()
		dlistItr = dlistItr.GetNext()
	}
}

func TestDoublyLinkedList_InsertAfter(t *testing.T) {
	dlist := NewLinkedList[int]()
	elements := dlist.AppendValue(1)
	_2n := dlist.InsertAfter(2, elements[0])
	_3n := dlist.InsertAfter(3, _2n)
	_4n := dlist.InsertAfter(4, _3n)
	dlist.InsertAfter(5, _4n)
	assert.Equal(t, int64(5), dlist.Len())

	dlist2 := list.New()
	_1n_2 := dlist2.PushBack(1)
	_2n_2 := dlist2.InsertAfter(2, _1n_2)
	_3n_2 := dlist2.InsertAfter(3, _2n_2)
	_4n_2 := dlist2.InsertAfter(4, _3n_2)
	dlist2.InsertAfter(5, _4n_2)

	assert.Equal(t, dlist.Len(), int64(dlist2.Len()))

	dlistItr := dlist.Front()
	dlist2Itr := dlist2.Front()
	for dlist2Itr != nil {
		assert.Equal(t, dlistItr.GetValue(), dlist2Itr.Value)
		dlist2Itr = dlist2Itr.Next()
		dlistItr = dlistItr.GetNext()
	}
}

func TestLinkedList_AppendValueThenRemove(t *testing.T) {
	t.Log("test linked list append value")
	dlist := NewLinkedList[int]()
	dlist2 := list.New()
	checkItems := func() {
		dlistItr := dlist.Front()
		dlist2Itr := dlist2.Front()
		for dlist2Itr != nil {
			assert.Equal(t, dlistItr.GetValue(), dlist2Itr.Value)
			dlist2Itr = dlist2Itr.Next()
			dlistItr = dlistItr.GetNext()
		}
	}

	elements := dlist.AppendValue(1, 2, 3, 4, 5)
	assert.Equal(t, len(elements), 5)
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, e: %v", idx, e)
		assert.Equal(t, elements[idx], e)
		t.Logf("addr: %p, return addr: %p", elements[idx], e)
	})

	dlist2.PushBack(1)
	dlist2.PushBack(2)
	_3n := dlist2.PushBack(3)
	dlist2.PushBack(4)
	dlist2.PushBack(5)
	assert.Equal(t, dlist.Len(), int64(dlist2.Len()))
	checkItems()

	t.Log("test linked list remove middle")
	dlist.Remove(elements[2])
	dlist2.Remove(_3n)
	checkItems()

	t.Log("test linked list remove head")
	dlist.Remove(dlist.Front())
	dlist2.Remove(dlist2.Front())
	checkItems()

	t.Log("test linked list remove tail")
	dlist.Remove(dlist.Back())
	dlist2.Remove(dlist2.Back())
	checkItems()

	t.Log("test linked list remove nil")
	dlist.Remove(nil)
	// dlist2.Remove(nil) // nil panic
	checkItems()

	t.Log("check released elements")
	assert.Equal(t, int64(dlist2.Len()), dlist.Len())
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, e: %v", idx, e)
	})
	for idx, e := range elements {
		t.Logf("index: %d, ptr: %p, e: %v", idx, e, e)
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/list
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkNewLinkedList_AppendValue
// BenchmarkNewLinkedList_AppendValue-4   	 3138679	       394.3 ns/op	      88 B/op	       3 allocs/op
func BenchmarkNewLinkedList_AppendValue(b *testing.B) {
	dlist := NewLinkedList[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dlist.AppendValue(i)
	}
	b.ReportAllocs()
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/list
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkSDKLinkedList_PushBack
// BenchmarkSDKLinkedList_PushBack-4   	 4632534	       237.2 ns/op	      55 B/op	       1 allocs/op
func BenchmarkSDKLinkedList_PushBack(b *testing.B) {
	dlist := list.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dlist.PushBack(i)
	}
	b.ReportAllocs()
}

func TestLinkedList_InsertAfterAndMove(t *testing.T) {
	dlist := NewLinkedList[int]()
	dlist2 := list.New()
	checkItems := func() {
		dlistItr := dlist.Front()
		dlist2Itr := dlist2.Front()
		assert.NotNil(t, dlistItr)
		assert.NotNil(t, dlist2Itr)
		assert.Equal(t, int64(dlist2.Len()), dlist.Len())
		for dlist2Itr != nil {
			assert.Equal(t, dlist2Itr.Value, dlistItr.GetValue())
			dlist2Itr = dlist2Itr.Next()
			dlistItr = dlistItr.GetNext()
		}
	}

	elements := dlist.AppendValue(1, 2, 3, 4, 5)
	_6n := dlist.InsertAfter(6, elements[len(elements)-1])
	_7n := dlist.InsertBefore(7, elements[0])
	assert.Equal(t, int64(7), dlist.Len())
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})

	dlist2.PushBack(1)
	dlist2.PushBack(2)
	dlist2.PushBack(3)
	dlist2.PushBack(4)
	dlist2.PushBack(5)
	_6n_2 := dlist2.InsertAfter(6, dlist2.Back())
	_7n_2 := dlist2.InsertBefore(7, dlist2.Front())
	assert.Equal(t, int64(dlist2.Len()), dlist.Len())
	checkItems()

	t.Log("test move after")
	dlist.MoveToBack(_7n)
	dlist2.MoveToBack(_7n_2)
	checkItems()
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})

	t.Log("test move to front")
	dlist.MoveToFront(_6n)
	dlist2.MoveToFront(_6n_2)
	checkItems()
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})

	t.Log("test move before")
	dlist.MoveBefore(_6n, _7n)
	dlist2.MoveBefore(_6n_2, _7n_2)
	checkItems()
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})

	t.Log("test move after")
	dlist.MoveAfter(_7n, dlist.Front())
	dlist2.MoveAfter(_7n_2, dlist2.Front())
	checkItems()
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})

	t.Log("test push front list")
	dlist_1 := NewLinkedList[int]()
	dlist_1.AppendValue(8, 9, 10)
	dlist2_1 := list.New()
	dlist2_1.PushBack(8)
	dlist2_1.PushBack(9)
	dlist2_1.PushBack(10)

	dlist.PushFrontList(dlist_1)
	dlist2.PushFrontList(dlist2_1)
	checkItems()
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})

	t.Log("test push back list")
	dlist_2 := NewLinkedList[int]()
	dlist_2.AppendValue(11, 12, 13)
	dlist2_2 := list.New()
	dlist2_2.PushBack(11)
	dlist2_2.PushBack(12)
	dlist2_2.PushBack(13)

	dlist.PushBackList(dlist_2)
	dlist2.PushBackList(dlist2_2)
	checkItems()
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})
}

func TestLinkedList_PushBack(t *testing.T) {
	dlist := NewLinkedList[int]()
	element := dlist.PushBack(1)
	assert.Equal(t, int64(1), dlist.Len())
	assert.Equal(t, element.GetValue(), 1)

	element = dlist.PushBack(2)
	assert.Equal(t, int64(2), dlist.Len())
	assert.Equal(t, element.GetValue(), 2)

	expected := []int{1, 2}
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		assert.Equal(t, expected[idx], e.GetValue())
	})

	reverseExpected := []int{2, 1}
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		assert.Equal(t, reverseExpected[idx], e.GetValue())
	})
}

func TestLinkedList_PushFront(t *testing.T) {
	dlist := NewLinkedList[int]()
	element := dlist.PushFront(1)
	assert.Equal(t, int64(1), dlist.Len())
	assert.Equal(t, element.GetValue(), 1)

	element = dlist.PushFront(2)
	assert.Equal(t, int64(2), dlist.Len())
	assert.Equal(t, element.GetValue(), 2)

	expected := []int{2, 1}
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		assert.Equal(t, expected[idx], e.GetValue())
	})

	reverseExpected := []int{1, 2}
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		assert.Equal(t, reverseExpected[idx], e.GetValue())
	})
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/list
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkDoublyLinkedList_PushBack
// BenchmarkDoublyLinkedList_PushBack-4   	 3944040	       258.9 ns/op	      64 B/op	       1 allocs/op
func BenchmarkDoublyLinkedList_PushBack(b *testing.B) {
	dlist := NewLinkedList[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dlist.PushBack(i)
	}
	b.ReportAllocs()
}

// goos: linux
// goarch: amd64
// pkg: github.com/benz9527/toy-box/algo/list
// cpu: Intel(R) Core(TM) i5-4590 CPU @ 3.30GHz
// BenchmarkDoublyLinkedList_Append
// BenchmarkDoublyLinkedList_Append-4   	 9450279	       127.1 ns/op	      16 B/op	       1 allocs/op
func BenchmarkDoublyLinkedList_Append(b *testing.B) {
	dlist := NewLinkedList[int]()
	elements := make([]NodeElement[int], 0, b.N)
	for i := 0; i < b.N; i++ {
		elements = append(elements, newNodeElement[int](i, dlist))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dlist.Append(elements[i])
	}
	b.StopTimer()
	b.ReportAllocs()
	assert.Equal(b, int64(b.N), dlist.Len())
}

func TestDoublyLinkedList_InsertBefore2(t *testing.T) {
	dlist := NewLinkedList[int]()
	_2n := dlist.InsertBefore(2, dlist.Front())
	assert.Equal(t, int64(1), dlist.Len())
	assert.Equal(t, _2n.GetValue(), 2)
	dlist.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("index: %d, addr: %p, e: %v", idx, e, e)
	})
	dlist.ReverseForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("reverse: index: %d, addr: %p, e: %v", idx, e, e)
	})
}
