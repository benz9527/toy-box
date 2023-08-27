package stack

// https://leetcode.cn/problems/trapping-rain-water/

func Trap(height []int) int {
	stack, ans := make([]int, 0, len(height)), 0
	peek := func() int {
		if len(stack) <= 0 {
			return -1
		}
		idx := stack[len(stack)-1]
		return idx
	}
	pop := func() int {
		if len(stack) == 0 {
			return -1
		}
		idx := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		return idx
	}
	push := func(idx int) {
		stack = append(stack, idx)
	}
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	push(0)
	for i := 1; i < len(height); i++ {
		for len(stack) > 0 && peek() != -1 && height[peek()] < height[i] {
			middle := pop() // 方便做横行计算
			if len(stack) > 0 && peek() != -1 {
				h := minimum(height[peek()], height[i]) - height[middle]
				w := i - peek() - 1
				ans += h * w
			}
		}
		push(i)
	}

	return ans
}
