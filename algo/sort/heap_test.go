package sort

import (
	"testing"
)

func TestHeapSort(t *testing.T) {
	DataComparator(HeapSort, 1000, 10, 100, true)
}
