package sort

import (
	"testing"
)

func TestShellSort(t *testing.T) {
	LogarithmicDetector(ShellSort, 1000, 10, 100, true)
}
