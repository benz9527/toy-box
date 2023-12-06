package sort

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	DataComparator(MergeSort, 1000, 100, 100, true)
}
