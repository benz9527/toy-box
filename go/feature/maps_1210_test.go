package feature

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Go1.21.0
// func clear[T ~[]Type | ~map[Type]Type1](t T)

type (
	tSlice[T any]             []T
	tMap[K comparable, V any] map[K]V
)

func TestGo12100MapClearFunc(t *testing.T) {
	arr := make([]tSlice[int], 0, 4)
	clear(arr)
	expected := make([]tSlice[int], 4)
	assert.NotEqual(t, expected, arr)
	assert.Equal(t, 0, len(arr))
	assert.Equal(t, 4, cap(arr))

	arr2 := make([]int, 0, 4)
	clear(arr2)
	expected2 := make([]int, 4)
	assert.NotEqual(t, expected2, arr2)
	assert.Equal(t, 0, len(arr2))
	assert.Equal(t, 4, cap(arr2))

	arr3 := []tSlice[int]{{1}, {2}, {3}, {4}}
	clear(arr3)
	expected3 := []tSlice[int]{{0}, {0}, {0}, {0}}
	assert.NotEqual(t, expected3, arr3)
	assert.Equal(t, 4, len(arr3))
	assert.Equal(t, 4, cap(arr3))

	m := tMap[string, int]{
		"abc":  1,
		"best": 2,
	}
	clear(m)
	assert.Equal(t, 0, len(m))
	assert.Equal(t, tMap[string, int]{}, m)
}
