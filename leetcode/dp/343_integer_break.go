package dp

// https://leetcode.cn/problems/integer-break/

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
