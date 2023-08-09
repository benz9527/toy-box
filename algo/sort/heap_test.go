package sort

import (
	"testing"
)

func TestHeapSort(t *testing.T) {
	LogarithmicDetector(HeapSort, 1000, 10, 100, true)
}
