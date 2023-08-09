package sort

import (
	"testing"
)

func TestBenchmarkSelectSort(t *testing.T) {
	LogarithmicDetector(SelectSort, 1000, 100, 100, true)
}
