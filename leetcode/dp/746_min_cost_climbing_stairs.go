package dp

// https://leetcode.cn/problems/min-cost-climbing-stairs/

func MinCostClimbingStairs(cost []int) int {
	n := len(cost)
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	dpFunc := func(start int) int {
		_dp := make([]int, n+1)
		for i := 0; i <= n; i++ {
			_dp[i] = -1
		}
		_dp[start] = 0
		for i := start; i < n; i++ {
			c := cost[i]
			if i <= n-1 {
				if _dp[i+1] < 0 {
					_dp[i+1] = _dp[i] + c
				} else {
					_dp[i+1] = minimum(_dp[i+1], _dp[i]+c)
				}
			}
			if i <= n-2 {
				if _dp[i+2] < 0 {
					_dp[i+2] = _dp[i] + c
				} else {
					_dp[i+2] = minimum(_dp[i+2], _dp[i]+c)
				}
			}
		}
		return _dp[n]
	}

	return minimum(dpFunc(0), dpFunc(1))
}

func MinCostClimbingStairsOptimize(cost []int) int {
	n := len(cost)
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	_dp := make([]int, n+1)
	for i := 2; i <= n; i++ {
		_dp[i] = minimum(_dp[i-1]+cost[i-1], _dp[i-2]+cost[i-2])
	}

	return _dp[n]
}
