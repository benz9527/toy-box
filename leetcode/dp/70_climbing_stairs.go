package dp

// https://leetcode.cn/problems/climbing-stairs

func ClimbStairs(n int) int {
	if n == 1 {
		return 1
	}
	_dp := make([]int, n+1)
	// value(_dp[i]) 是登上该阶楼梯的总方案数量
	// 假设不知道 value(_dp[i]) 的具体值，只能往下试图获取
	// 登上 i-1 阶楼梯的方案数，以此为基础再加新方案
	// 根据题目的意思，一次可以登一阶或者两阶
	// 则 value(_dp[i]) = value(_dp[i-1]+1)
	// 或者 value(_dp[i]) = value(_dp[i-2]+2)
	// 这样的逻辑说明，加方案的来源就是固定的，就两个
	// value(_dp[i]) = value(_dp[i-1]) + value(_dp[i-2])
	// n >= 1，说明没有 0 阶方案
	_dp[1], _dp[2] = 1, 2
	for i := 3; i <= n; i++ {
		_dp[i] = _dp[i-1] + _dp[i-2]
	}
	return _dp[n]
}

func ClimbStairsOptimize(n int) int {
	if n == 1 {
		return 1
	}
	_dp := []int{1: 1, 2: 2}
	for i := 3; i <= n; i++ {
		sum := _dp[1] + _dp[2]
		_dp[1], _dp[2] = _dp[2], sum
	}
	return _dp[2]
}

func ClimbStairsOptimizeByDp(n int) int {
	if n == 1 {
		return 1
	}
	_dp := make([]int, n+1)
	_dp[1] = 1
	_dp[2] = 2
	for j := 3; j <= n; j++ {
		for i := 1; i <= 2; i++ {
			_dp[j] += _dp[j-i]
		}
	}
	return _dp[n]
}
