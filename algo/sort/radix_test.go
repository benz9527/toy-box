package sort

import (
	"testing"
)

func TestRadixSort(t *testing.T) {
	DataComparator(RadixSort, 1000, 10, 500, true)
}
