package sort

import (
	"testing"
)

func TestShellSort(t *testing.T) {
	DataComparator(ShellSort, 1000, 10, 100, true)
}
