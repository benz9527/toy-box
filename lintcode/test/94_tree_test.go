package test

import (
	"github.com/benz9527/toy-box/lintcode/tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxPathSum_3(t *testing.T) {
	type TreeNode = tree.Node94
	tr := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
		},
		Right: &TreeNode{
			Val: 3,
		},
	}
	res := tree.MaxPathSum(tr)
	assert.Equal(t, 6, res)
}

func TestMaxPathSum_negative_3(t *testing.T) {
	type TreeNode = tree.Node94
	tr := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: -2,
		},
		Right: &TreeNode{
			Val: 3,
		},
	}
	res := tree.MaxPathSum(tr)
	assert.Equal(t, 4, res)
}

func TestMaxPathSum_hyper_7(t *testing.T) {
	type TreeNode = tree.Node94
	tr := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 4,
			},
		},
		Right: &TreeNode{
			Val: -5,
			Left: &TreeNode{
				Val: 5,
			},
			Right: &TreeNode{
				Val: 6,
			},
		},
	}
	res := tree.MaxPathSum(tr)
	assert.Equal(t, 8, res)
}

func TestMaxPathSum_7(t *testing.T) {
	type TreeNode = tree.Node94
	tr := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 4,
			},
			Right: &TreeNode{
				Val: 5,
			},
		},
		Right: &TreeNode{
			Val: 3,
			Left: &TreeNode{
				Val: 6,
			},
			Right: &TreeNode{
				Val: 7,
			},
		},
	}
	res := tree.MaxPathSum(tr)
	assert.Equal(t, 18, res)
}

func TestMaxPathSum_negative_7(t *testing.T) {
	type TreeNode = tree.Node94
	tr := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 100,
			Left: &TreeNode{
				Val: -8,
			},
			Right: &TreeNode{
				Val: -99,
			},
		},
		Right: &TreeNode{
			Val: 5,
			Left: &TreeNode{
				Val: 6,
			},
			Right: &TreeNode{
				Val: -7,
			},
		},
	}
	res := tree.MaxPathSum(tr)
	assert.Equal(t, 112, res)
}

func TestMaxPathSum_all_negative(t *testing.T) {
	type TreeNode = tree.Node94
	tr := &TreeNode{
		Val: -10,
		Left: &TreeNode{
			Val: -20,
			Right: &TreeNode{
				Val: -31,
				Left: &TreeNode{
					Val: -24,
				},
				Right: &TreeNode{
					Val: -5,
					Left: &TreeNode{
						Val: -6,
						Left: &TreeNode{
							Val: -8,
						},
						Right: &TreeNode{
							Val: -9,
						},
					},
					Right: &TreeNode{
						Val: -7,
					},
				},
			},
		},
	}
	res := tree.MaxPathSum(tr)
	assert.Equal(t, -5, res)
}
