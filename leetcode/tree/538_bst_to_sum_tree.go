package tree

type Node struct {
	Val         int
	Left, Right *Node
}

func ConvertBST(root *Node) *Node {
	traverse(root)
	return root
}

var (
	sum = 0
)

func traverse(root *Node) {
	if root == nil {
		return
	}
	traverse(root.Right)
	root.Val += sum
	sum = root.Val
	traverse(root.Left)
}
