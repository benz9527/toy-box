package sort

import (
	"testing"
)

func TestInsertSortBenchmark(t *testing.T) {
	DataComparator(InsertSort, 1000, 1000, 100, true)
}

func TestInsertSortBenchmarkSimplify(t *testing.T) {
	DataComparator(InsertSortSimplify, 1000, 10, 100, true)
}
