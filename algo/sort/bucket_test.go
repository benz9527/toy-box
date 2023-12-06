package sort

import (
	"testing"
)

func TestBucketSort(t *testing.T) {
	DataComparator(BucketSort, 1000, 10, 100, true)
}
