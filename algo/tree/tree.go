package tree

import "fmt"

type Node struct {
	Item        int
	Height      int
	Left, Right *Node
}

type BinaryTree struct {
	Root *Node
}

func CreateBinaryTree() *BinaryTree {
	return &BinaryTree{Root: nil}
}

/*
Recursive implementation
*/

func PostOrder(n *Node) {
	if n == nil {
		return
	}

	PostOrder(n.Left)
	PostOrder(n.Right)
	fmt.Print(n.Item, " -> ")
}

func InOrder(n *Node) {
	if n == nil {
		return
	}

	InOrder(n.Left)
	fmt.Print(n.Item, " -> ")
	InOrder(n.Right)
}

func PreOrder(n *Node) {
	if n == nil {
		return
	}

	fmt.Print(n.Item, " -> ")
	PreOrder(n.Left)
	PreOrder(n.Right)
}

/*
loop implementation by stack (non-recursive)
*/

/*
General loop implementation (non-recursive)
*/
