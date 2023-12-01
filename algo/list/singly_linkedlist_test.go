package list_test

import (
	"testing"
	"time"

	"github.com/benz9527/toy-box/algo/list"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestSinglyLinkedListSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Singly Linked List Suite",
		types.SuiteConfig{
			LabelFilter:     "SinglyLinkedList",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

var _ = ginkgo.Describe("Singly Linked List Unit Tests", ginkgo.Ordered, ginkgo.Serial, func() {
	ginkgo.It("singly linked list generation",
		ginkgo.Label("SinglyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			slist.InsertValue(1)
			slist.InsertValue(2)
			slist.InsertValue(3)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			slist.ForEach(func(e list.NodeElement[int]) {
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
		},
	)
	ginkgo.It("singly linked list generation with insert before/after",
		ginkgo.Label("SinglyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.InsertValue(1)
			_2n := slist.InsertBefore(_1n, 2)
			slist.InsertAfter(_2n, 3)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			expected := []int{2, 3, 1}
			actual := make([]int, 0, 3)
			slist.ForEach(func(e list.NodeElement[int]) {
				actual = append(actual, e.GetValue())
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
			assert.Equal(ginkgo.GinkgoT(), expected, actual)
		},
	)
	ginkgo.It("singly linked list generation then find target node",
		ginkgo.Label("SinglyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.InsertValue(1)
			_2n := slist.InsertBefore(_1n, 2)
			slist.InsertAfter(_2n, 3)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			_3n, ok := slist.Find(3)
			assert.True(ginkgo.GinkgoT(), ok)
			assert.Equal(ginkgo.GinkgoT(), 3, _3n.GetValue())
		},
	)
	ginkgo.It("singly linked list generation then remove target node",
		ginkgo.Label("SinglyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.InsertValue(1)
			_2n := slist.InsertBefore(_1n, 2)
			_3n := slist.InsertAfter(_2n, 3)
			slist.Insert(list.NewSinglyNodeElement(4))
			slist.Remove(_3n)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			expected := []int{2, 1, 4}
			actual := make([]int, 0, 3)
			slist.ForEach(func(e list.NodeElement[int]) {
				actual = append(actual, e.GetValue())
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
			assert.Equal(ginkgo.GinkgoT(), expected, actual)
		},
	)
	ginkgo.It("singly linked list generation, remove target node and insert nil node",
		ginkgo.Label("SinglyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.InsertValue(1)
			_2n := slist.InsertBefore(_1n, 2)
			_3n := slist.InsertAfter(_2n, 3)
			slist.Insert(list.NewSinglyNodeElement(4))
			slist.Remove(_3n)
			slist.Insert(nil)
			slist.InsertBefore(nil, 5)
			slist.InsertAfter(nil, 6)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			expected := []int{2, 1, 4}
			actual := make([]int, 0, 3)
			slist.ForEach(func(e list.NodeElement[int]) {
				actual = append(actual, e.GetValue())
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
			assert.Equal(ginkgo.GinkgoT(), expected, actual)
		},
	)
})
