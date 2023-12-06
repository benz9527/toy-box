package sort

import (
	"testing"
)

func TestBenchmarkSelectSort(t *testing.T) {
	DataComparator(SelectSort, 1000, 100, 100, true)
}
