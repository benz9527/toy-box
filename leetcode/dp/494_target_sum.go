package dp

// https://leetcode.cn/problems/target-sum/

// 给你一个非负整数数组 nums 和一个整数 target 。
// 向数组中的每个整数前添加 '+' 或 '-' ，然后串联起所有整数，可以构造一个 表达式 ：
// 例如，nums = [2, 1] ，可以在 2 之前添加 '+' ，在 1 之前添加 '-' ，然后串联起来得到表达式 "+2-1" 。
// 返回可以通过上述方法构造的、运算结果等于 target 的不同 表达式 的数目。

// 这里总和 sum 是可知的，因为都是正整数的元素，直接求和，sum = nums[0] + ... + nums[n]
// 假设计算 sum 的左加数（元素和）为 left，右加数（元素和）为 right，而且 left 是较大的一方
// 则 sum = left + right
// 如果 sum > target（正整数），才有机会对 target 进行计算，而且 target = left - right
// 以上是正整数的推断
// 如果 target 是负数，我们需要先判断 sum 的总数是否能小于 target
// sum = (-nums[0]) + ... + (-nums[n])
// 假设计算 sum 的左加数（元素和）为 left，右加数（元素和）为 right，两者都是负数，而且 left 是较小的一方
// 则 sum = left + right
// 如果 sum < target (负整数)，才有机会对 target 进行计算，而且 target = left - right
// 转换为绝对值处理话，abs(sum) > abs(target)，则 abs(target) = abs(target) - abs(right)
// 这样就说明 target 不管正负，都是一样处理，可以统一逻辑:
// abs(left) - (abs(sum) - abs(left)) = abs(target)
// abs(left) = (abs(sum) + abs(target)) / 2
// 这样就转换为求左表达式的可能取值

func FindTargetSumWays(nums []int, target int) int {
	sum := 0
	for _, n := range nums {
		sum += n // abs(sum)
	}

	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}
	target = abs(target)
	if target > sum {
		return 0 // 没法计算
	}
	if (target+sum)&0x1 == 1 {
		// 不管是向上还是向下取整，数目都对不上
		return 0
	}

	left := (target + sum) / 2
	_dp := make([]int, left+1)
	_dp[0] = 1
	for _, n := range nums {
		for j := left; j >= n; j-- {
			_dp[j] += _dp[j-n]
		}
	}
	return _dp[left]
}

// 回溯法暴力搜索

func FindTargetSumWaysByBackTracing(nums []int, target int) int {
	ways := 0
	// 0: +, 1: -
	calc := make([]int, len(nums))
	stack := make([]int, 0, len(nums))
	pop := func() int { // 一定返回索引
		idx := len(stack) - 1
		if idx < 0 {
			return -1
		}
		stack = stack[0:idx]
		calc[idx]++
		return idx
	}
	push := func(item int) {
		stack = append(stack, item) // 一定是记录索引
	}
	sum := func(items ...int) int {
		res := 0
		for _, i := range items {
			res += i
		}
		return res
	}

	for i := 0; i < len(nums); i++ {
		push(nums[i])
	}

	if target == sum(stack...) {
		ways++
	}
	for idx := pop(); idx != -1; idx = pop() {
		if calc[idx] > 0 && calc[idx]&0x1 == 0 { // 转向策略
			continue
		}
		for len(stack) < len(nums) {
			item := nums[len(stack)] // 当前栈的长度就是要 push 进去元素的索引
			if calc[len(stack)]&0x1 == 1 {
				item = -item
			}
			push(item)
		}
		if target == sum(stack...) { // 最后执行操作关键逻辑
			ways++
		}
	}
	return ways
}
