package container_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"

	"github.com/benz9527/toy-box/toys/pkg/container"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
)

func TestZeroCopyFile(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Thread Safe Map Suite",
		types.SuiteConfig{
			LabelFilter:     "ThreadSafeMap",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

type testCloser struct {
	closeByErr bool
	id         string
	io.Writer
}

func (t *testCloser) Close() error {
	_, _ = t.Writer.Write(bytes.NewBufferString(t.id + " closed\n").Bytes())
	if t.closeByErr {
		return io.ErrUnexpectedEOF
	}
	return nil
}

var _ = ginkgo.Describe("Thread Safe Map Element io.Closer",
	ginkgo.Label("ThreadSafeMap"),
	func() {
		ginkgo.It("should close all elements", func() {
			m := container.NewThreadSafeMap[*testCloser]()
			m.AddOrUpdate("1", &testCloser{id: "1", Writer: ginkgo.GinkgoWriter})
			m.AddOrUpdate("2", &testCloser{id: "2", closeByErr: true, Writer: ginkgo.GinkgoWriter})
			err := m.Close()
			assert.NoError(ginkgo.GinkgoT(), err)
		})
	},
)
