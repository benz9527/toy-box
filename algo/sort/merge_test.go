package sort

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	LogarithmicDetector(MergeSort, 1000, 100, 100, true)
}
