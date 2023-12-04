package list_test

import (
	clist "container/list"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/benz9527/toy-box/algo/list"
)

func TestLinkedList_AppendValue(t *testing.T) {
	dlist := list.NewLinkedList[int]()
	elements := dlist.AppendValue(1, 2, 3, 4, 5)
	assert.Equal(t, len(elements), 5)
	dlist.ForEach(func(idx int64, e list.NodeElement[int]) {
		t.Logf("index: %d, e: %v", idx, e)
		assert.Equal(t, elements[idx], e)
		t.Logf("addr: %p, return addr: %p", elements[idx], e)
	})

	dlist2 := clist.New()
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
	dlist := list.NewLinkedList[int]()
	elements := dlist.AppendValue(1)
	_2n := dlist.InsertBefore(2, elements[0])
	_3n := dlist.InsertBefore(3, _2n)
	_4n := dlist.InsertBefore(4, _3n)
	dlist.InsertBefore(5, _4n)
	assert.Equal(t, int64(5), dlist.Len())

	dlist2 := clist.New()
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
	dlist := list.NewLinkedList[int]()
	elements := dlist.AppendValue(1)
	_2n := dlist.InsertAfter(2, elements[0])
	_3n := dlist.InsertAfter(3, _2n)
	_4n := dlist.InsertAfter(4, _3n)
	dlist.InsertAfter(5, _4n)
	assert.Equal(t, int64(5), dlist.Len())

	dlist2 := clist.New()
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
	dlist := list.NewLinkedList[int]()
	dlist2 := clist.New()
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
	dlist.ForEach(func(idx int64, e list.NodeElement[int]) {
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
	dlist.ForEach(func(idx int64, e list.NodeElement[int]) {
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
	dlist := list.NewLinkedList[int]()
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
	dlist := clist.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dlist.PushBack(i)
	}
	b.ReportAllocs()
}
