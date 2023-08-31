package api

import (
	"github.com/stretchr/testify/assert"
	isort "sort"
	"strconv"
	"strings"
	"testing"
)

func TestStringToCharArrAndSort(t *testing.T) {
	src := "8901211113"
	arr := strings.Split(src, "")
	isort.Strings(arr)
	src = strings.Join(arr, "")
	// 连续数值 + 滑动窗口
	other := "8910111213"
	arr2 := strings.Split(other, "")
	isort.Strings(arr2)
	other = strings.Join(arr2, "")
	assert.Equal(t, src, other)
}

func TestStringToNumber(t *testing.T) {
	snumber := "100"
	n, err := strconv.Atoi(snumber)
	assert.NoError(t, err)
	assert.Equal(t, 100, n)
}
