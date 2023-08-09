package test

import (
	_sort "github.com/benz9527/toy-box/algo/sort"
	"github.com/benz9527/toy-box/lintcode/sort"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortIntegers(t *testing.T) {
	src := []int{3, 2, 1, 4, 5}
	expected := []int{1, 2, 3, 4, 5}
	src1 := _sort.CopyArr(src, len(src))
	sort.SortIntegers(src1)
	assert.Equal(t, expected, src1)

	src1 = _sort.CopyArr(src, len(src))
	sort.SortIntegersByAlgoBubble(src1)
	assert.Equal(t, expected, src1)
}
