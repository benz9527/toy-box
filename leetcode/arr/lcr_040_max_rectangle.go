package arr

// https://leetcode.cn/problems/PLYXKQ/
// 困难

func MaximalRectangle(matrix []string) int {
	maxArea := 0
	if len(matrix) <= 0 {
		return maxArea
	}
	for i := len(matrix) - 1; i >= 0; i-- {
		heights := make([]int, len(matrix[0]))
		for j := 0; j < len(matrix[0]); j++ {
			h := i
			// 不断地缩减高度（换底），获取以该层为底构建高度序列
			for h >= 0 && matrix[h][j] == '1' {
				// 出现 0 的说明高度不连续
				h--
				heights[j]++
			}
		}
		maxArea = maximum(MaxRecArea(heights), maxArea)
	}
	return maxArea
}

func maximum(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxRecArea(heights []int) int {
	maxArea := 0
	// 栈操作
	stack := make([]int, 0, len(heights))
	peek := func() int {
		return stack[len(stack)-1]
	}
	push := func(item int) {
		stack = append(stack, item)
	}
	pop := func() int {
		res := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		return res
	}

	// 遍历高度
	for i := 0; i < len(heights); i++ {
		// 进栈前的过滤，如果将要进入的元素小于栈顶元素对应的高度值
		// 则会破坏自底向上的不断递增的要求，导致弹栈
		for len(stack) > 0 && heights[i] <= heights[peek()] {
			j := pop()
			k := -1 // -1 表示到达最左边，栈底
			if len(stack) > 0 {
				k = peek()
			}
			// i-1 就是当前元素的前一位是之前高度连通域的右边界
			// 每一个右边界都会计算一次，直到左边界
			curArea := (i - k - 1) * heights[j]
			maxArea = maximum(curArea, maxArea)
		}
		// 默认每个高度的索引作为元素都要进栈
		push(i)
	}
	// 由于存在最后没有小的元素入栈，导致栈内元素滞留，也需要计算滞留的
	// 连通域大小
	for len(stack) > 0 {
		j := pop()
		k := -1
		if len(stack) > 0 {
			k = peek()
		}
		// len(heights) - 1 就是已经到达整体的最右边界
		curArea := (len(heights) - k - 1) * heights[j]
		maxArea = maximum(curArea, maxArea)
	}
	return maxArea
}
