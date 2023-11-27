package concurrency_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConcurrency(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Concurrency Suite")
}
