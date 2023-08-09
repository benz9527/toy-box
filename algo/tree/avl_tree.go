package tree

func CreateAVLNode(item int) *Node {
	return &Node{Item: item, Height: 1, Left: nil, Right: nil}
}

func AVLHeight(n *Node) int {
	if n == nil {
		return 0
	}
	return n.Height
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func LeftRotate(xn *Node) *Node {
	yn := xn.Right
	temp := yn.Left

	xn.Right = temp
	yn.Left = xn

	xn.Height = Max(AVLHeight(xn.Left), AVLHeight(xn.Right)) + 1
	yn.Height = Max(AVLHeight(yn.Left), AVLHeight(yn.Right)) + 1
	return yn
}

func RightRotate(yn *Node) *Node {
	xn := yn.Left
	temp := xn.Right

	yn.Left = temp
	xn.Right = yn

	xn.Height = Max(AVLHeight(xn.Left), AVLHeight(xn.Right)) + 1
	yn.Height = Max(AVLHeight(yn.Left), AVLHeight(yn.Right)) + 1
	return xn
}

func GetBalanceFactor(n *Node) int {
	if n == nil {
		return 0
	}
	return AVLHeight(n.Left) - AVLHeight(n.Right)
}

func AVLInsertNode(n *Node, newItem int) *Node {
	if n == nil {
		return CreateAVLNode(newItem)
	}

	if newItem < n.Item {
		n.Left = AVLInsertNode(n.Left, newItem)
	} else if newItem > n.Item {
		n.Right = AVLInsertNode(n.Right, newItem)
	} else {
		return n
	}

	n.Height = 1 + Max(AVLHeight(n.Left), AVLHeight(n.Right))
	bf := GetBalanceFactor(n)
	if bf > 1 {
		if newItem < n.Left.Item {
			return RightRotate(n)
		} else if newItem > n.Left.Item {
			n.Left = LeftRotate(n.Left)
			return RightRotate(n)
		}
	}

	if bf < -1 {
		if newItem > n.Right.Item {
			return LeftRotate(n)
		} else if newItem < n.Right.Item {
			n.Right = RightRotate(n.Right)
			return LeftRotate(n)
		}
	}
	return n
}

func AVLMinimumValueNode(n *Node) *Node {
	cur := n
	for cur.Left != nil {
		cur = cur.Left
	}
	return cur
}

func AVLDeleteNode(n *Node, deletedItem int) *Node {
	if n == nil {
		return n
	}

	if deletedItem < n.Item {
		n.Left = AVLDeleteNode(n.Left, deletedItem)
	} else if deletedItem > n.Item {
		n.Right = AVLDeleteNode(n.Right, deletedItem)
	} else {
		// n will be removed, n is a leaf.
		if n.Left == nil ||
			n.Right == nil {
			var temp *Node = nil

			if temp == n.Left {
				temp = n.Right
			} else {
				temp = n.Left
			}

			if temp == nil {
				temp = n
				n = nil
			} else {
				n = temp
			}
		} else {
			temp := AVLMinimumValueNode(n.Right)
			n.Item = temp.Item
			n.Right = AVLDeleteNode(n.Right, temp.Item)
		}
	}

	if n == nil {
		return n
	}

	n.Height = Max(AVLHeight(n.Left), AVLHeight(n.Right)) + 1
	bf := GetBalanceFactor(n)
	if bf > 1 {
		if GetBalanceFactor(n.Left) >= 0 {
			return RightRotate(n)
		} else {
			n.Left = LeftRotate(n.Left)
			return RightRotate(n)
		}
	}

	if bf < -1 {
		if GetBalanceFactor(n.Right) <= 0 {
			return LeftRotate(n)
		} else {
			n.Right = RightRotate(n.Right)
			return LeftRotate(n)
		}
	}
	return n
}
