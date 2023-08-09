package tree

func CreateNodeWithItem(item int) *Node {
	return &Node{Item: item}
}

func FullBinaryTreeNodeValidate(n *Node) bool {
	if n == nil {
		return true
	}

	if n.Left == nil && n.Right == nil {
		return true
	}

	if n.Left != nil && n.Right != nil {
		return FullBinaryTreeNodeValidate(n.Left) && FullBinaryTreeNodeValidate(n.Right)
	}

	return false
}

func IsFullBinaryTree(bt *BinaryTree) bool {
	if bt == nil {
		return true
	}

	return FullBinaryTreeNodeValidate(bt.Root)
}
