package dp

// https://leetcode.cn/problems/maximum-subarray/
// 连续的子序列（隐含直接有序，不能插空）

// Time exceeded

func MaxSumOfSubArray(nums []int) int {
	cols := len(nums)
	rows := cols

	res := -10000
	for i := 0; i < cols; i++ {
		lastdp := nums[i]
		for j := i; j < rows; j++ {
			if i != j {
				lastdp += nums[j]
			}
			if lastdp > res {
				res = lastdp
			}
		}
	}

	return res
}

func MaxSumOfSubArrayOptimize(nums []int) int {
	cols := len(nums)
	_dp := make([]int, cols)
	_dp[0] = nums[0]
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	res := _dp[0]
	for i := 1; i < cols; i++ {
		_dp[i] = maximum(_dp[i-1]+nums[i], nums[i])
		if res < _dp[i] {
			res = _dp[i]
		}
	}

	return res
}
