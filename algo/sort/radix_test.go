package sort

import (
	"testing"
)

func TestRadixSort(t *testing.T) {
	LogarithmicDetector(RadixSort, 1000, 10, 500, true)
}
