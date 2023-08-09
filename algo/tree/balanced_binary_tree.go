package tree

func IsHeightBalanced(n *Node, h int) (bool, int) {
	if n == nil {
		h = 0
		return true, 0
	}

	l, lh := IsHeightBalanced(n.Left, 0)
	r, rh := IsHeightBalanced(n.Right, 0)

	if lh > rh {
		h = lh + 1
	} else {
		h = rh + 1
	}

	if lh-rh >= 2 ||
		rh-lh >= 2 {
		return false, h
	}

	return l && r, h
}
