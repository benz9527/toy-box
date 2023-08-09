package sort

import (
	"testing"
)

func TestQuickSort(t *testing.T) {
	LogarithmicDetector(QuickSort, 1000, 1000, 100, true)
}

func TestQuickSort2(t *testing.T) {
	LogarithmicDetector(QuickSort2, 1000, 1000, 100, true)
}
