package tree

import "fmt"

const RED int = 1
const BLACK int = 0

type RedBlackNode struct {
	Item                int
	Color               int
	Parent, Left, Right *RedBlackNode
}

// EmptyNil All value is default nil(zero) value
var EmptyNil *RedBlackNode = nil

type RedBlackTree struct {
	Root *RedBlackNode
}

func (rbTree *RedBlackTree) PreOrder(n *RedBlackNode) {
	if n != EmptyNil {
		fmt.Print(n.Item, " ")
		rbTree.PreOrder(n.Left)
		rbTree.PreOrder(n.Right)
	}
}

func (rbTree *RedBlackTree) InOrder(n *RedBlackNode) {
	if n != EmptyNil {
		rbTree.InOrder(n.Left)
		fmt.Print(n.Item, " ")
		rbTree.InOrder(n.Right)
	}
}

func (rbTree *RedBlackTree) PostOrder(n *RedBlackNode) {
	if n != EmptyNil {
		rbTree.PostOrder(n.Left)
		rbTree.PostOrder(n.Right)
		fmt.Print(n.Item, " ")
	}
}

func (rbTree *RedBlackTree) Search(n *RedBlackNode, k int) *RedBlackNode {
	if n == EmptyNil || k == n.Item {
		return n
	}

	if k < n.Item {
		return rbTree.Search(n.Left, k)
	}
	return rbTree.Search(n.Right, k)
}

func (rbTree *RedBlackTree) LeftRotate(nx *RedBlackNode) {
	ny := nx.Right
	nx.Right = ny.Left

	if ny.Left != EmptyNil {
		ny.Left.Parent = nx
	}

	ny.Parent = nx.Parent
	if nx.Parent == nil {
		rbTree.Root = ny
	} else if nx == nx.Parent.Left {
		nx.Parent.Left = ny
	} else {
		nx.Parent.Right = ny
	}

	ny.Left = nx
	nx.Parent = ny
}

func (rbTree *RedBlackTree) RightRotate(nx *RedBlackNode) {
	ny := nx.Left
	nx.Left = ny.Right

	if ny.Right != EmptyNil {
		ny.Right.Parent = nx
	}

	ny.Parent = nx.Parent
	if nx.Parent == nil {
		rbTree.Root = ny
	} else if nx == nx.Parent.Right {
		nx.Parent.Right = ny
	} else {
		nx.Parent.Left = ny
	}
	ny.Right = nx
	nx.Parent = ny
}

func (rbTree *RedBlackTree) Insert(k int) {
	newNode := &RedBlackNode{
		Parent: nil,
		Item:   k,
		Left:   EmptyNil,
		Right:  EmptyNil,
		Color:  RED,
	}

	var ny *RedBlackNode = nil
	var nx = rbTree.Root

	for nx != EmptyNil {
		ny = nx
		if newNode.Item < nx.Item {
			nx = nx.Left
		} else {
			nx = nx.Right
		}
	}

	newNode.Parent = ny
	if ny == nil {
		rbTree.Root = newNode
	} else if newNode.Item < ny.Item {
		ny.Left = newNode
	} else {
		ny.Right = newNode
	}

	if newNode.Parent == nil {
		newNode.Color = BLACK
		return
	}

	if newNode.Parent.Parent == nil {
		return
	}

	rbTree.PostInsertBalance(newNode)
}

// PostInsertBalance Balance the node after insertion
func (rbTree *RedBlackTree) PostInsertBalance(target *RedBlackNode) {
	var tmp *RedBlackNode = nil
	for target.Parent.Color == RED {
		if target.Parent == target.Parent.Parent.Right {
			tmp = target.Parent.Parent.Left
			if tmp.Color == RED {
				tmp.Color = BLACK
				target.Parent.Color = BLACK
				target.Parent.Parent.Color = RED
				target = target.Parent.Parent
			} else {
				if target == target.Parent.Left {
					target = target.Parent
					rbTree.RightRotate(target)
				}
				target.Parent.Color = BLACK
				target.Parent.Parent.Color = RED
				rbTree.LeftRotate(target.Parent.Parent)
			}
		} else {
			tmp = target.Parent.Parent.Right

			if tmp.Color == RED {
				tmp.Color = BLACK
				target.Parent.Color = BLACK
				target.Parent.Parent.Color = RED
				target = target.Parent.Parent
			} else {
				if target == target.Parent.Right {
					target = target.Parent
					rbTree.LeftRotate(target)
				}

				target.Parent.Color = BLACK
				target.Parent.Parent.Color = RED
				rbTree.RightRotate(target.Parent.Parent)
			}
		}

		if target == rbTree.Root {
			break
		}
	}
	rbTree.Root.Color = BLACK
}

func (rbTree *RedBlackTree) Transplant(currentNode, childNode *RedBlackNode) {
	if currentNode.Parent == nil {
		rbTree.Root = childNode
	} else if currentNode == currentNode.Parent.Left {
		currentNode.Parent.Left = childNode
	} else {
		currentNode.Parent.Right = childNode
	}
	childNode.Parent = currentNode.Parent
}

func (rbTree *RedBlackTree) Minimum(n *RedBlackNode) *RedBlackNode {
	for n.Left != EmptyNil {
		n = n.Left
	}
	return n
}

func (rbTree *RedBlackTree) Maximum(n *RedBlackNode) *RedBlackNode {
	for n.Right != EmptyNil {
		n = n.Right
	}
	return n
}

func (rbTree *RedBlackTree) Successor(nx *RedBlackNode) *RedBlackNode {
	if nx.Right != EmptyNil {
		return rbTree.Minimum(nx.Right)
	}

	ny := nx.Parent
	for ny != EmptyNil &&
		nx == ny.Right {
		nx = ny
		ny = ny.Parent
	}
	return ny
}

func (rbTree *RedBlackTree) Predecessor(nx *RedBlackNode) *RedBlackNode {
	if nx.Left != EmptyNil {
		return rbTree.Maximum(nx.Left)
	}

	ny := nx.Parent
	for ny != EmptyNil &&
		nx == ny.Left {
		nx = ny
		ny = ny.Parent
	}
	return ny
}

func (rbTree *RedBlackTree) Delete(parentNode *RedBlackNode, k int) {
	nz := EmptyNil
	var nx, ny *RedBlackNode = nil, nil
	for parentNode != EmptyNil {
		if parentNode.Item == k {
			nz = parentNode
		}

		if parentNode.Item <= k {
			parentNode = parentNode.Right
		} else {
			parentNode = parentNode.Left
		}
	}

	if nz == EmptyNil {
		fmt.Println("Could not find key in the tree!")
		return
	}

	ny = nz
	nyOldColor := ny.Color
	if nz.Left == EmptyNil {
		nx = nz.Right
		rbTree.Transplant(nz, nz.Right)
	} else if nz.Right == EmptyNil {
		nx = nz.Left
		rbTree.Transplant(nz, nz.Left)
	} else {
		ny = rbTree.Minimum(nz.Right)
		nyOldColor = ny.Color
		nx = ny.Right
		if ny.Parent == nz {
			nx.Parent = ny
		} else {
			rbTree.Transplant(ny, ny.Right)
			ny.Right = nz.Right
			ny.Right.Parent = ny
		}

		rbTree.Transplant(nz, ny)
		ny.Left = nz.Left
		ny.Left.Parent = ny
		ny.Color = nz.Color
	}

	if nyOldColor == BLACK {
		rbTree.PostDeleteBalance(nx)
	}
}

// PostDeleteBalance Balance the tree after deletion of a node
func (rbTree *RedBlackTree) PostDeleteBalance(nx *RedBlackNode) {
	var ns *RedBlackNode = nil
	for nx != rbTree.Root &&
		nx.Color == BLACK {
		if nx == nx.Parent.Left {
			ns = nx.Parent.Right
			if ns.Color == RED {
				ns.Color = BLACK
				nx.Parent.Color = RED
				rbTree.LeftRotate(nx.Parent)
				ns = nx.Parent.Right
			}

			if ns.Left.Color == BLACK &&
				ns.Right.Color == BLACK {
				ns.Color = RED
				nx = nx.Parent
			} else {
				if ns.Right.Color == BLACK {
					ns.Left.Color = BLACK
					ns.Color = RED
					rbTree.RightRotate(ns)
					ns = nx.Parent.Right
				}

				ns.Color = nx.Parent.Color
				nx.Parent.Color = BLACK
				ns.Right.Color = BLACK
				rbTree.LeftRotate(nx.Parent)
				nx = rbTree.Root
			}
		} else {
			ns = nx.Parent.Left
			if ns.Color == RED {
				ns.Color = BLACK
				nx.Parent.Color = RED
				rbTree.RightRotate(nx.Parent)
				ns = nx.Parent.Left
			}

			if ns.Left.Color == BLACK &&
				ns.Right.Color == BLACK {
				ns.Color = RED
				nx = nx.Parent
			} else {
				if ns.Left.Color == BLACK {
					ns.Right.Color = BLACK
					ns.Color = RED
					rbTree.LeftRotate(ns)
					ns = nx.Parent.Left
				}

				ns.Color = nx.Parent.Color
				nx.Parent.Color = BLACK
				ns.Left.Color = BLACK
				rbTree.RightRotate(nx.Parent)
				nx = rbTree.Root
			}
		}
	}
	nx.Color = BLACK
}
