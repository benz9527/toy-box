package tree

func Depth(n *Node) int {
	d := 0
	for n != nil {
		d++
		n = n.Left
	}
	return d
}

func isPerfectBinaryTree(n *Node, depth int, level int) bool {
	if n == nil {
		return true
	}

	if n.Left == nil && n.Right == nil {
		return depth == (level + 1)
	}

	if n.Left == nil || n.Right == nil {
		return false
	}

	// Not a tail recursion.
	return isPerfectBinaryTree(n.Left, depth, level+1) &&
		isPerfectBinaryTree(n.Right, depth, level+1)
}

func IsPerfectBinaryTree(n *Node) bool {
	depth := Depth(n)
	return isPerfectBinaryTree(n, depth, 0)
}
