package dp

// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-with-cooldown/

func MaxProfitV(prices []int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([][]int, len(prices))
	for i := 0; i < len(prices); i++ {
		_dp[i] = make([]int, 4)
	}
	_dp[0][0] = -prices[0]
	for i := 1; i < len(prices); i++ {
		// 状态机转移图
		_dp[i][0] = maximum(_dp[i-1][0], maximum(_dp[i-1][2], _dp[i-1][3])-prices[i])
		_dp[i][1] = maximum(_dp[i-1][1], _dp[i-1][0]+prices[i])
		_dp[i][2] = _dp[i-1][1]
		_dp[i][3] = maximum(_dp[i-1][3], _dp[i-1][2])
	}

	return maximum(_dp[len(prices)-1][1],
		maximum(_dp[len(prices)-1][2], _dp[len(prices)-1][3]))
}

func MaxProfitVOptimize(prices []int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([]int, 4)
	_dp[0] = -prices[0]
	for i := 1; i < len(prices); i++ {
		// 状态机转移图
		buyMax, beforeSaleMax, freezeMax, keepSaleMax := _dp[0], _dp[1], _dp[2], _dp[3]
		_dp[0] = maximum(buyMax, maximum(freezeMax, keepSaleMax)-prices[i])
		_dp[1] = maximum(beforeSaleMax, buyMax+prices[i])
		_dp[2] = beforeSaleMax
		_dp[3] = maximum(keepSaleMax, freezeMax)
	}

	return maximum(_dp[1], maximum(_dp[2], _dp[3]))
}
