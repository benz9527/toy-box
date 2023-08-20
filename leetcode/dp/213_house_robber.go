package dp

// https://leetcode.cn/problems/house-robber-ii/

// 你是一个专业的小偷，计划偷窃沿街的房屋，每间房内都藏有一定的现金。这个地方所有的房屋都 围成一圈 ，这意味着第一个房屋和最后一个房屋是紧挨着的。同时，相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警 。
// 给定一个代表每个房屋存放金额的非负整数数组，计算你 在不触动警报装置的情况下 ，今晚能够偷窃到的最高金额。
// 1 <= nums.length <= 100
// 0 <= nums[i] <= 1000

func RobII(nums []int) int {
	houses := len(nums)
	if houses == 1 {
		return nums[0]
	}
	if houses == 2 {
		return 0
	}
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	robII := func(numList []int) int {
		h := len(numList)
		_dp := make([]int, h)

		// _dp[1] = max(numList[0], numList[1])
		// _dp[2] = max(_dp[0]+numList[2], _dp[1])
		_dp[0] = numList[0]
		if h > 1 {
			_dp[1] = maximum(numList[0], numList[1])
		}
		for i := 2; i < h; i++ {
			for j := h - 1; j >= i; j-- {
				_dp[j] = maximum(_dp[j-2]+numList[j], _dp[j-1])
			}
		}

		return _dp[h-1]
	}
	// 去掉一个数据，靠算两遍来补回去
	// 成环的情况：其实只有一种大范围就是首尾都取不到
	// 在此基础上分裂成两种小情况来修正：
	// 1. 不取头
	// 2. 不取尾
	res1 := robII(nums[0 : len(nums)-1])
	res2 := robII(nums[1:])
	return maximum(res1, res2)
}
