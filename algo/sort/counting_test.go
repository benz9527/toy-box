package sort

import (
	"testing"
)

func TestCountingSort(t *testing.T) {
	DataComparator(CountingSort, 1000, 10, 100, true)
}
