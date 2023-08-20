package dp

// https://leetcode.cn/problems/house-robber/

// 经典题目
// 间隔选取物品，和间隔爬楼梯类似的做法

func Rob(nums []int) int {
	houses := len(nums)
	_dp := make([]int, houses)
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// _dp[2] = _dp[1] or _dp[2] = _dp[0] + nums[2]
	// _dp[1] = nums[1] or nums[0]
	_dp[0] = nums[0]
	if houses > 1 {
		_dp[1] = maximum(nums[0], nums[1])
	}

	for i := 2; i < houses; i++ {
		for j := houses - 1; j >= i; j-- {
			_dp[j] = maximum(_dp[j-1], _dp[j-2]+nums[j])
		}
	}
	return _dp[houses-1]
}
