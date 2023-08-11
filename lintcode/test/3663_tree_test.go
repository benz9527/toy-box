package test

import (
	"github.com/benz9527/toy-box/lintcode/tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUndirectedTreeDiameter(t *testing.T) {
	edges := [][]int{
		{0, 1}, {1, 2}, {2, 3}, {1, 4}, {4, 5},
	}
	res := tree.UndirectedTreeDiameter(edges)
	assert.Equal(t, 4, res)
}

func TestUndirectedTreeDiameter2(t *testing.T) {
	edges := [][]int{
		{11, 3}, {5, 11}, {14, 5}, {6, 5}, {13, 14}, {4, 14}, {10, 4}, {1, 11}, {9, 13}, {8, 14}, {7, 5}, {0, 4}, {2, 7}, {12, 14},
	}
	res := tree.UndirectedTreeDiameter(edges)
	assert.Equal(t, 5, res)
}

func TestUndirectedTreeDiameter4(t *testing.T) {
	edges := [][]int{
		{3, 17}, {19, 3}, {10, 19}, {1, 10}, {9, 1}, {22, 1}, {2, 9}, {26, 19}, {16, 9}, {6, 2}, {15, 22},
		{20, 16}, {18, 2}, {0, 17}, {7, 18}, {8, 20}, {23, 6}, {27, 17}, {24, 15}, {25, 15}, {12, 15},
		{29, 7}, {4, 6}, {13, 4}, {14, 27}, {21, 23}, {28, 0}, {11, 26}, {5, 1},
	}
	res := tree.UndirectedTreeDiameter(edges)
	assert.Equal(t, 11, res)
}

func TestUndirectedTreeDiameter5(t *testing.T) {
	edges := [][]int{
		{0, 1}, {1, 2},
	}
	res := tree.UndirectedTreeDiameter(edges)
	assert.Equal(t, 2, res)
}
