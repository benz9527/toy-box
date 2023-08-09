package tree

type Node94 struct {
	Val   int
	Left  *Node94
	Right *Node94
}

func MaxPathSum(root *Node94) int {
	if root == nil {
		return 0
	}
	var (
		pathValues = make([]int, 0, 128)
		lv, rv     int
		lok, rok   bool
	)
	lv, pathValues, lok = maxSingleSidePathSum(root.Left, pathValues)
	rv, pathValues, rok = maxSingleSidePathSum(root.Right, pathValues)
	sumList := make([]int, 0, 8)
	if lok {
		sumList = append(sumList, lv, root.Val+lv)
	}
	if rok {
		sumList = append(sumList, rv, root.Val+rv)
	}
	if lok && rok {
		sumList = append(sumList, root.Val+lv+rv)
	}
	pathValues = append(pathValues, sumList...)
	res := root.Val
	for _, v := range pathValues {
		if v > res {
			res = v
		}
	}
	return res
}

func maxSingleSidePathSum(root *Node94, pathValues []int) (int, []int, bool) {
	if root == nil {
		return 0, pathValues, false
	}
	if root.Left == nil && root.Right == nil {
		pathValues = append(pathValues, root.Val)
		return root.Val, pathValues, true
	}

	// post order travel
	var (
		lv, rv   int
		lok, rok bool
	)
	lv, pathValues, lok = maxSingleSidePathSum(root.Left, pathValues)
	rv, pathValues, rok = maxSingleSidePathSum(root.Right, pathValues)
	res := root.Val
	sumList := make([]int, 0, 8)
	sumList = append(sumList, res)
	if lok {
		if lv >= 0 {
			sumList = append(sumList, root.Val+lv)
		}
		if res < 0 && lv < 0 && res < lv {
			sumList = append(sumList, lv)
		}
	}
	if rok {
		if rv >= 0 {
			sumList = append(sumList, root.Val+rv)
		}
		if res < 0 && rv < 0 && res < rv {
			sumList = append(sumList, rv)
		}
	}
	if lok && rok {
		pathValues = append(pathValues, root.Val+lv+rv)
	}
	for _, v := range sumList {
		if v > res {
			res = v
		}
	}
	pathValues = append(pathValues, sumList...)
	return res, pathValues, true
}
