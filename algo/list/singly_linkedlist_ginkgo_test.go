package list_test

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/benz9527/toy-box/algo/list"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestSinglyLinkedListSuite(t *testing.T) {
	// FIXME ginkgo unable run parallel specs in the same package
	type testCase struct {
		name string
		suit func()
	}
	testcases := []testCase{
		{
			name: "1",
			suit: func() {
				gomega.RegisterFailHandler(ginkgo.Fail)
				ginkgo.RunSpecs(t, "Singly Linked BasicLinkedList Suite",
					types.SuiteConfig{
						LabelFilter:     "singlyLinkedList",
						ParallelTotal:   1,
						ParallelProcess: 1,
						GracePeriod:     5 * time.Second,
					},
					types.ReporterConfig{
						Verbose: true,
					},
				)
			},
		},
		{
			name: "2",
			suit: func() {
				gomega.RegisterFailHandler(ginkgo.Fail)
				ginkgo.RunSpecs(t, "Concurrent Singly Linked BasicLinkedList Suite",
					types.SuiteConfig{
						LabelFilter:     "ConcurrentSinglyLinkedList Parallel",
						ParallelTotal:   1,
						ParallelProcess: 1,
						GracePeriod:     5 * time.Second,
					},
					types.ReporterConfig{
						Verbose: true,
					},
				)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tt.suit()
		})
	}
}

var _ = ginkgo.Describe("Singly Linked BasicLinkedList Unit Tests", ginkgo.Ordered, ginkgo.Serial, func() {
	ginkgo.It("singly linked list generation",
		ginkgo.Label("singlyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			slist.AppendValue(1, 2, 3)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			slist.ForEach(func(idx int64, e list.NodeElement[int]) {
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
		},
	)
	ginkgo.It("singly linked list generation with insert before/after",
		ginkgo.Label("singlyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.AppendValue(1)
			_2n := slist.InsertBefore(2, _1n[0])
			slist.InsertAfter(3, _2n)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			expected := []int{2, 3, 1}
			actual := make([]int, 0, 3)
			slist.ForEach(func(idx int64, e list.NodeElement[int]) {
				actual = append(actual, e.GetValue())
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
			assert.Equal(ginkgo.GinkgoT(), expected, actual)
		},
	)
	ginkgo.It("singly linked list generation then find target node",
		ginkgo.Label("singlyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.AppendValue(1)
			_2n := slist.InsertBefore(2, _1n[0])
			slist.InsertAfter(3, _2n)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			_3n, ok := slist.FindFirst(3)
			assert.True(ginkgo.GinkgoT(), ok)
			assert.Equal(ginkgo.GinkgoT(), 3, _3n.GetValue())
		},
	)
	ginkgo.It("singly linked list generation then remove target node",
		ginkgo.Label("singlyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.AppendValue(1)
			_2n := slist.InsertBefore(2, _1n[0])
			_3n := slist.InsertAfter(3, _2n)
			slist.Append(list.NewNodeElement(4))
			_3n_1 := slist.Remove(_3n)
			assert.Equal(ginkgo.GinkgoT(), _3n.GetValue(), _3n_1.GetValue())
			_3n_1 = slist.Remove(list.NewNodeElement(3))
			assert.Nil(ginkgo.GinkgoT(), _3n_1)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			expected := []int{2, 1, 4}
			actual := make([]int, 0, 3)
			slist.ForEach(func(idx int64, e list.NodeElement[int]) {
				actual = append(actual, e.GetValue())
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
			assert.Equal(ginkgo.GinkgoT(), expected, actual)
		},
	)
	ginkgo.It("singly linked list generation, remove target node and insert nil node",
		ginkgo.Label("singlyLinkedList"),
		func() {
			slist := list.NewSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			_1n := slist.AppendValue(1)
			_2n := slist.InsertBefore(2, _1n[0])
			_3n := slist.InsertAfter(3, _2n)
			slist.Append(list.NewNodeElement(4))
			slist.Remove(_3n)
			slist.Append(nil)
			slist.InsertBefore(5, nil)
			slist.InsertAfter(6, nil)
			e := slist.Remove(nil)
			assert.Nil(ginkgo.GinkgoT(), e)
			assert.Equal(ginkgo.GinkgoT(), int64(3), slist.Len())
			expected := []int{2, 1, 4}
			actual := make([]int, 0, 3)
			slist.ForEach(func(idx int64, e list.NodeElement[int]) {
				actual = append(actual, e.GetValue())
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
			assert.Equal(ginkgo.GinkgoT(), expected, actual)
		},
	)
})

var _ = ginkgo.Describe("Concurrent Singly Linked BasicLinkedList Unit Tests", func() {
	ginkgo.It("singly linked list generation, remove target node and run in parallel",
		ginkgo.Label("ConcurrentSinglyLinkedList Parallel"),
		func() {
			slist := list.NewConcurrentSinglyLinkedList[int]()
			assert.NotNil(ginkgo.GinkgoT(), slist)
			wg := sync.WaitGroup{}
			wg.Add(5)
			go func() {
				slist.AppendValue(1)
				wg.Done()
			}()
			go func() {
				slist.AppendValue(2)
				wg.Done()
			}()
			go func() {
				slist.AppendValue(3)
				wg.Done()
			}()
			go func() {
				_3n, ok := slist.FindFirst(3)
				if ok && _3n != nil {
					slist.Remove(_3n)
				}
				wg.Done()
			}()
			go func() {
				slist.AppendValue(4)
				wg.Done()
			}()
			wg.Wait()
			expected1 := []int{1, 2, 4}
			expected2 := []int{1, 2, 3, 4}
			actual := make([]int, 0, 3)
			slist.ForEach(func(idx int64, e list.NodeElement[int]) {
				actual = append(actual, e.GetValue())
				ginkgo.GinkgoWriter.Printf("node element: %v\n", e.GetValue())
			})
			sort.Ints(actual)
			if slist.Len() == 3 {
				assert.Equal(ginkgo.GinkgoT(), expected1, actual)
			} else if slist.Len() == 4 {
				assert.Equal(ginkgo.GinkgoT(), expected2, actual)
			} else {
				assert.Fail(ginkgo.GinkgoT(), "unexpected singly linked list length")
			}
		},
	)
})
