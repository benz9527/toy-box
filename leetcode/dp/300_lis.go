package dp

func LengthOfLIS(nums []int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	n := len(nums)
	_dp := make([]int, n)
	// 每个背包最短的最长子序列长度是 1
	for i := 0; i < n; i++ {
		_dp[i] = 1
	}
	// permutation problem
	res := 1                 // 全局最长
	for i := 1; i < n; i++ { // knapsacks
		for j := 0; j < i; j++ { // sub knapsacks
			// 静态背包作为标杆去收集动态更小背包的结果
			// 相当于两次循环都在遍历背包
			if nums[i] > nums[j] {
				_dp[i] = maximum(_dp[i], _dp[j]+1)
			}
		}
		if _dp[i] > res {
			res = _dp[i]
		}
	}
	return res
}
