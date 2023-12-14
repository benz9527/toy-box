package test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
	"testing"
	"time"
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
