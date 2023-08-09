package sort

import (
	"testing"
)

func TestBucketSort(t *testing.T) {
	LogarithmicDetector(BucketSort, 1000, 10, 100, true)
}
