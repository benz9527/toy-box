package test

import (
	"github.com/benz9527/toy-box/lintcode/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEarliestAcq(t *testing.T) {
	src := [][]int{
		{20220101, 0, 1}, {20220109, 3, 4}, {20220304, 0, 4}, {20220305, 0, 3}, {20220404, 2, 4},
	}
	size := 6
	res := graph.EarliestAcq(src, size)
	assert.Equal(t, 20220404, res)

	src = [][]int{
		{7, 3, 1}, {3, 0, 3}, {2, 0, 1}, {1, 1, 2}, {5, 3, 2},
	}
	size = 4
	res = graph.EarliestAcq(src, size)
	assert.Equal(t, 3, res)

	src = [][]int{
		{5040, 0, 0}, {4569, 0, 1}, {5719, 0, 2}, {6628, 0, 3}, {4446, 0, 4}, {2991, 0, 5}, {6649, 0, 6}, {9064, 0, 7}, {2439, 0, 8}, {7674, 1, 0}, {9431, 1, 1}, {9240, 1, 2}, {2979, 1, 3}, {7823, 1, 4}, {993, 1, 5}, {4555, 1, 6}, {6100, 1, 7}, {2598, 1, 8}, {719, 2, 0}, {7901, 2, 1}, {7062, 2, 2}, {5657, 2, 3}, {9388, 2, 4}, {4388, 2, 5}, {5367, 2, 6}, {5952, 2, 7}, {5416, 2, 8}, {3879, 3, 0}, {1285, 3, 1}, {530, 3, 2}, {6719, 3, 3}, {6307, 3, 4}, {7702, 3, 5}, {2907, 3, 6}, {7987, 3, 7}, {4645, 3, 8}, {5648, 4, 0}, {7151, 4, 1}, {9413, 4, 2}, {6413, 4, 3}, {37, 4, 4}, {9720, 4, 5}, {3772, 4, 6}, {8165, 4, 7}, {1095, 4, 8}, {4765, 5, 0}, {1907, 5, 1}, {2435, 5, 2}, {6142, 5, 3}, {4827, 5, 4}, {2548, 5, 5}, {8363, 5, 6}, {1237, 5, 7}, {8762, 5, 8}, {3535, 6, 0}, {2782, 6, 1}, {5447, 6, 2}, {9027, 6, 3}, {1711, 6, 4}, {6633, 6, 5}, {3626, 6, 6}, {3852, 6, 7}, {691, 6, 8}, {8211, 7, 0}, {5674, 7, 1}, {1908, 7, 2}, {4560, 7, 3}, {120, 7, 4}, {6239, 7, 5}, {5367, 7, 6}, {1715, 7, 7}, {247, 7, 8}, {7970, 8, 0}, {2680, 8, 1}, {1719, 8, 2}, {754, 8, 3}, {466, 8, 4}, {9679, 8, 5}, {2226, 8, 6}, {1074, 8, 7}, {8365, 8, 8}, {2497, 9, 0}, {9950, 9, 1}, {3398, 9, 2}, {9403, 9, 3}, {2724, 9, 4}, {6987, 9, 5}, {9552, 9, 6}, {6933, 9, 7}, {1167, 9, 8},
	}
	size = 10
	res = graph.EarliestAcq(src, size)
	assert.Equal(t, 1237, res)

	src = [][]int{
		{10, 0, 1},
	}
	size = 2
	res = graph.EarliestAcq(src, size)
	assert.Equal(t, 10, res)
}

func TestEarliestAcq2(t *testing.T) {
	src := [][]int{
		{20220101, 0, 1}, {20220109, 3, 4}, {20220304, 0, 4}, {20220305, 0, 3}, {20220404, 2, 4},
	}
	size := 6
	res := graph.EarliestAcq2(src, size)
	assert.Equal(t, 20220404, res)

	src = [][]int{
		{7, 3, 1}, {3, 0, 3}, {2, 0, 1}, {1, 1, 2}, {5, 3, 2},
	}
	size = 4
	res = graph.EarliestAcq2(src, size)
	assert.Equal(t, 3, res)

	src = [][]int{
		{5040, 0, 0}, {4569, 0, 1}, {5719, 0, 2}, {6628, 0, 3}, {4446, 0, 4}, {2991, 0, 5}, {6649, 0, 6}, {9064, 0, 7}, {2439, 0, 8}, {7674, 1, 0}, {9431, 1, 1}, {9240, 1, 2}, {2979, 1, 3}, {7823, 1, 4}, {993, 1, 5}, {4555, 1, 6}, {6100, 1, 7}, {2598, 1, 8}, {719, 2, 0}, {7901, 2, 1}, {7062, 2, 2}, {5657, 2, 3}, {9388, 2, 4}, {4388, 2, 5}, {5367, 2, 6}, {5952, 2, 7}, {5416, 2, 8}, {3879, 3, 0}, {1285, 3, 1}, {530, 3, 2}, {6719, 3, 3}, {6307, 3, 4}, {7702, 3, 5}, {2907, 3, 6}, {7987, 3, 7}, {4645, 3, 8}, {5648, 4, 0}, {7151, 4, 1}, {9413, 4, 2}, {6413, 4, 3}, {37, 4, 4}, {9720, 4, 5}, {3772, 4, 6}, {8165, 4, 7}, {1095, 4, 8}, {4765, 5, 0}, {1907, 5, 1}, {2435, 5, 2}, {6142, 5, 3}, {4827, 5, 4}, {2548, 5, 5}, {8363, 5, 6}, {1237, 5, 7}, {8762, 5, 8}, {3535, 6, 0}, {2782, 6, 1}, {5447, 6, 2}, {9027, 6, 3}, {1711, 6, 4}, {6633, 6, 5}, {3626, 6, 6}, {3852, 6, 7}, {691, 6, 8}, {8211, 7, 0}, {5674, 7, 1}, {1908, 7, 2}, {4560, 7, 3}, {120, 7, 4}, {6239, 7, 5}, {5367, 7, 6}, {1715, 7, 7}, {247, 7, 8}, {7970, 8, 0}, {2680, 8, 1}, {1719, 8, 2}, {754, 8, 3}, {466, 8, 4}, {9679, 8, 5}, {2226, 8, 6}, {1074, 8, 7}, {8365, 8, 8}, {2497, 9, 0}, {9950, 9, 1}, {3398, 9, 2}, {9403, 9, 3}, {2724, 9, 4}, {6987, 9, 5}, {9552, 9, 6}, {6933, 9, 7}, {1167, 9, 8},
	}
	size = 10
	res = graph.EarliestAcq2(src, size)
	assert.Equal(t, 1237, res)

	src = [][]int{
		{10, 0, 1},
	}
	size = 2
	res = graph.EarliestAcq2(src, size)
	assert.Equal(t, 10, res)
}
