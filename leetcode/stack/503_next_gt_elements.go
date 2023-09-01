package stack

// https://leetcode.cn/problems/next-greater-element-ii/description/

// 给定一个循环数组 nums （ nums[nums.length - 1] 的下一个元素是 nums[0] ），返回 nums 中每个元素的 下一个更大元素 。
// 数字 x 的 下一个更大的元素 是按数组遍历顺序，这个数字之后的第一个比它更大的数，这意味着你应该循环地搜索它的下一个更大的数。
// 如果不存在，则输出 -1 。

func NextGreaterElements(nums []int) []int {
	stack := make([]int, 0, len(nums))
	pop := func() int {
		// index
		idx := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		return idx
	}
	push := func(n int) {
		stack = append(stack, n)
	}
	isEmpty := func() bool {
		return len(stack) == 0
	}
	peek := func() int {
		if isEmpty() {
			return -1
		}
		return stack[len(stack)-1]
	}

	ans := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		ans[i] = -1
	}
	n := len(nums)
	for i := 0; i < n*2; i++ {
		for !isEmpty() && peek() != -1 && nums[peek()%n] < nums[i%n] {
			idx := pop()
			ans[idx%n] = nums[i%n]
		}
		push(i)
	}
	return ans
}
