package tree

import (
	"fmt"
	"testing"
)

func generateBinaryTree() *BinaryTree {
	bt := CreateBinaryTree()
	bt.Root = &Node{Item: 1}
	bt.Root.Left = &Node{Item: 12}
	bt.Root.Right = &Node{Item: 9}
	bt.Root.Left.Left = &Node{Item: 5}
	bt.Root.Left.Right = &Node{Item: 6}
	return bt
}

func TestBinaryTreePostOrder(t *testing.T) {
	bt := generateBinaryTree()
	PostOrder(bt.Root)
}

func TestBinaryTreeInOrder(t *testing.T) {
	bt := generateBinaryTree()
	InOrder(bt.Root)
}

func TestBinaryTreePreOrder(t *testing.T) {
	bt := generateBinaryTree()
	PreOrder(bt.Root)
}

func TestFullBinaryTreeValidation(t *testing.T) {
	fbt := CreateBinaryTree()
	fbt.Root = CreateNodeWithItem(1)
	fbt.Root.Left = CreateNodeWithItem(2)
	fbt.Root.Right = CreateNodeWithItem(3)
	fbt.Root.Left.Left = CreateNodeWithItem(4)
	fbt.Root.Left.Right = CreateNodeWithItem(5)
	fbt.Root.Right.Left = CreateNodeWithItem(6)
	fbt.Root.Right.Right = CreateNodeWithItem(7)

	fmt.Println(IsFullBinaryTree(fbt))
}

func TestIsHeightBalanced(t *testing.T) {
	bt := CreateBinaryTree()
	bt.Root = CreateNodeWithItem(1)
	bt.Root.Left = CreateNodeWithItem(2)
	bt.Root.Right = CreateNodeWithItem(3)
	bt.Root.Left.Left = CreateNodeWithItem(4)
	bt.Root.Left.Right = CreateNodeWithItem(5)
	bt.Root.Left.Left.Left = CreateNodeWithItem(6)
	fmt.Println(IsHeightBalanced(bt.Root, 0))
}

func TestAVLTree(t *testing.T) {
	avl := CreateBinaryTree()
	avl.Root = AVLInsertNode(avl.Root, 33)
	avl.Root = AVLInsertNode(avl.Root, 13)
	avl.Root = AVLInsertNode(avl.Root, 53)
	avl.Root = AVLInsertNode(avl.Root, 9)
	avl.Root = AVLInsertNode(avl.Root, 21)
	avl.Root = AVLInsertNode(avl.Root, 61)
	avl.Root = AVLInsertNode(avl.Root, 8)
	avl.Root = AVLInsertNode(avl.Root, 11)

	PreOrder(avl.Root)
	fmt.Println()
	avl.Root = AVLDeleteNode(avl.Root, 13)
	PreOrder(avl.Root)
	fmt.Println()
}

func TestAVLTree2(t *testing.T) {
	avl := CreateBinaryTree()
	avl.Root = AVLInsertNode(avl.Root, 33)
	avl.Root = AVLInsertNode(avl.Root, 21)
	avl.Root = AVLInsertNode(avl.Root, 53)
	avl.Root = AVLInsertNode(avl.Root, 8)
	avl.Root = AVLInsertNode(avl.Root, 29)
	avl.Root = AVLInsertNode(avl.Root, 49)
	avl.Root = AVLInsertNode(avl.Root, 61)
	avl.Root = AVLInsertNode(avl.Root, 27)
	avl.Root = AVLInsertNode(avl.Root, 30)
	avl.Root = AVLInsertNode(avl.Root, 40)
	avl.Root = AVLInsertNode(avl.Root, 5)
	avl.Root = AVLInsertNode(avl.Root, 23)

	PreOrder(avl.Root)
	fmt.Println()
	avl.Root = AVLDeleteNode(avl.Root, 5)
	PreOrder(avl.Root)
	fmt.Println()
}

func TestBTree(t *testing.T) {
	btree := NewBTree(2)
	btree.Insert(8)
	btree.Insert(9)
	btree.Insert(10)
	btree.Insert(11)
	btree.Insert(15)
	btree.Insert(20)
	btree.Insert(17)

	btree.Show()

	if btree.Contain(12) {
		fmt.Println("found")
	} else {
		fmt.Println("not found")
	}
}

func TestRemoveKey(t *testing.T) {
	btree := NewBTree(3)
	btree.Insert(1)
	btree.Insert(3)
	btree.Insert(7)
	btree.Insert(10)
	btree.Insert(11)
	btree.Insert(13)
	btree.Insert(14)
	btree.Insert(15)
	btree.Insert(18)
	btree.Insert(16)
	btree.Insert(19)
	btree.Insert(24)
	btree.Insert(25)
	btree.Insert(26)
	btree.Insert(21)
	btree.Insert(4)
	btree.Insert(5)
	btree.Insert(20)
	btree.Insert(22)
	btree.Insert(2)
	btree.Insert(17)
	btree.Insert(12)
	btree.Insert(6)

	btree.PreorderShow()

	if btree.Contain(12) {
		fmt.Println("found")
	} else {
		fmt.Println("not found")
	}

	fmt.Printf("\nremove 6\n")
	btree.RemoveKey(6)
	btree.PreorderShow()

	fmt.Printf("\nremove 13\n")
	btree.RemoveKey(13)
	btree.PreorderShow()
}
