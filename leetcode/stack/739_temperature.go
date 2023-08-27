package stack

// https://leetcode.cn/problems/daily-temperatures
// 单调栈，只会从右边找符合关系的元素，栈内元素不符合要求即刻弹出（栈存的索引）

func DailyTemperatures(temperatures []int) []int {
	stack, ans := make([]int, 0, len(temperatures)), make([]int, len(temperatures))
	peek := func() int {
		if len(stack) <= 0 {
			return -1
		}
		idx := stack[len(stack)-1]
		return idx
	}
	pop := func() int {
		idx := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		return idx
	}
	push := func(idx int) {
		stack = append(stack, idx)
	}

	for i := 0; i < len(temperatures); i++ {
		if len(stack) == 0 || temperatures[peek()] >= temperatures[i] {
			push(i)
			continue
		}
		for peek() != -1 && temperatures[peek()] < temperatures[i] {
			idx := pop()
			ans[idx] = i - idx
		}
		push(i)
	}

	return ans
}
