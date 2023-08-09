package sort

import (
	"testing"
)

func TestCountingSort(t *testing.T) {
	LogarithmicDetector(CountingSort, 1000, 10, 100, true)
}
