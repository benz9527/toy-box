package dp

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// BFS 这个无法处理退化为链表的树，导致

func RobIII(root *TreeNode) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	getSubNodes := func(nodes ...*TreeNode) []*TreeNode {
		subNodes := make([]*TreeNode, 0, 2*len(nodes))
		for _, n := range nodes {
			if n.Left != nil {
				subNodes = append(subNodes, n.Left)
			}
			if n.Right != nil {
				subNodes = append(subNodes, n.Right)
			}
		}
		return subNodes
	}
	rootHouses := make([]int, 0, 16)
	rootLessHouses := make([]int, 0, 16)
	bfs := func(r *TreeNode) {
		if r == nil {
			return
		}
		rootHouses = append(rootHouses, r.Val)
		i := 0
		for nodes := getSubNodes(r); len(nodes) > 0; nodes = getSubNodes(nodes...) {
			i++
			if i&0x1 == 1 {
				for _, n := range nodes {
					rootLessHouses = append(rootLessHouses, n.Val)
				}
			} else {
				for _, n := range nodes {
					rootHouses = append(rootHouses, n.Val)
				}
			}
		}
	}
	bfs(root)
	rob := func(nums []int) int {
		houses := len(nums)
		mx := 0
		for i := 0; i < houses; i++ {
			mx += nums[i]
		}
		return mx
	}
	res1 := rob(rootHouses)
	res2 := rob(rootLessHouses)
	return maximum(res1, res2)
}

// 树形动态规划，在树上进行递归公式的推导。

func RobIIIByPostOrder(root *TreeNode) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	var postOrder func(*TreeNode) (int, int)
	postOrder = func(node *TreeNode) (int, int) {
		// 递归退出条件
		if node == nil {
			// 第一个返回值是偷取当前结点时能取得的最大值
			// 第二个返回值是不偷取当前结点时能取得的最大值
			return 0, 0
		}
		ls, lns := postOrder(node.Left)
		rs, rns := postOrder(node.Right)
		// 偷取当前结点时，必须获取左右子树不被偷取时的值
		// 不偷取当前结点时，左右子树可偷可不偷，只要取最大即可
		return node.Val + lns + rns, maximum(ls, lns) + maximum(rs, rns)
	}
	s, ns := postOrder(root)
	return maximum(s, ns)
}
