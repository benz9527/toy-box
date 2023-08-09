package tree

func CountNumNodes(n *Node) int {
	if n == nil {
		return 0
	}

	return 1 + CountNumNodes(n.Left) + CountNumNodes(n.Right)
}

func IsCompleteBinaryTree(n *Node, idx int, num int) bool {
	if n == nil {
		return true
	}

	if idx >= num {
		return false
	}

	return IsCompleteBinaryTree(n.Left, 2*idx+1, num) &&
		IsCompleteBinaryTree(n.Right, 2*idx+2, num)
}
