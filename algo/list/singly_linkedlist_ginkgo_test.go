package list

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSinglyLinkedList(t *testing.T) {
	l1 := NewSinglyLinkedList[int]()
	assert.NotNil(t, l1)
	elements := l1.AppendValue(1, 2, 3)
	assert.Equal(t, int64(3), l1.Len())
	l1.ForEach(func(idx int64, e NodeElement[int]) {
		t.Logf("node element: %v\n", e.GetValue())
	})

	l2 := list.New()
	l2.PushBack(1)
	_l2_e2 := l2.PushBack(2)
	l2.PushBack(3)
	assert.Equal(t, int64(3), int64(l2.Len()))

	l1Itr := l1.(*singlyLinkedList[int]).getRootHead()
	l2Itr := l2.Front()
	for l2Itr != nil {
		assert.Equal(t, l1Itr.GetValue(), l2Itr.Value)
		l2Itr = l2Itr.Next()
		l1Itr = l1Itr.GetNext()
	}

	l1.Remove(elements[1])
	l2.Remove(_l2_e2)
	assert.Equal(t, int64(2), l1.Len())
	l1Itr = l1.(*singlyLinkedList[int]).getRootHead()
	l2Itr = l2.Front()
	for l2Itr != nil {
		assert.Equal(t, l1Itr.GetValue(), l2Itr.Value)
		l2Itr = l2Itr.Next()
		l1Itr = l1Itr.GetNext()
	}

}
