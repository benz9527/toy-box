package dp

// https://leetcode.cn/problems/integer-break/

// 给定一个正整数 n ，将其拆分为 k 个 正整数 的和（ k >= 2 ），并使这些整数的乘积最大化。

func IntegerBreak(n int) int {
	_dp := make([]int, n+1)
	_dp[2] = 1
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for i := 2; i <= n; i++ {
		for j := i; j <= n; j++ {
			_dp[j] = maximum(_dp[j], i*_dp[j-i])
			_dp[j] = maximum(_dp[j], i*(j-i))
		}
	}
	return _dp[n]
}

func IntegerBreakOptimize(n int) int {
	_dp := make([]int, n+1)
	_dp[2] = 1
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for i := 3; i <= n; i++ {
		for j := 1; j <= i/2; j++ {
			_dp[i] = maximum(_dp[i], maximum(j*_dp[i-j], j*(i-j)))
		}
	}
	return _dp[n]
}
