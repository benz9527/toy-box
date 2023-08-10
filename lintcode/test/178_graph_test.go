package test

import (
	"github.com/benz9527/toy-box/lintcode/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipleSlice(t *testing.T) {
	edges := [][]int{
		{1, 2}, {3, 4}, {5, 6},
	}
	t.Log(len(edges))
	t.Log(len(edges[0]))
}

func TestValidTree(t *testing.T) {
	edges := [][]int{
		{1, 0},
	}
	res := graph.ValidTree(2, edges)
	assert.Equal(t, true, res)

	edges = [][]int{
		{1, 0},
	}
	res = graph.ValidTree(3, edges)
	assert.Equal(t, false, res)

	edges = [][]int{
		{1, 0}, {1, 2},
		{3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 3},
	}
	res = graph.ValidTree(8, edges)
	assert.Equal(t, false, res)

	edges = [][]int{
		{0, 1}, {1, 2}, {3, 2}, {4, 3}, {4, 5}, {5, 6}, {6, 7},
	}
	res = graph.ValidTree(8, edges)
	assert.Equal(t, true, res)
}

func TestValidTreeCompare(t *testing.T) {
	edges := [][]int{
		{1, 0}, {3, 0}, {3, 2}, {1, 2},
	}
	res := graph.ValidTree(4, edges)
	assert.Equal(t, false, res)
	res = graph.ValidTree3(4, edges)
	assert.Equal(t, false, res)

	edges = [][]int{
		{1, 0}, {3, 0}, {1, 2},
	}
	res = graph.ValidTree(4, edges)
	assert.Equal(t, true, res)
	res = graph.ValidTree3(4, edges)
	assert.Equal(t, true, res)
}
