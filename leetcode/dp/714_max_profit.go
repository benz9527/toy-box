package dp

// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-with-transaction-fee/

func MaxProfitVI(prices []int, fee int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([]int, 2)
	_dp[0] = -prices[0]
	for i := 1; i < len(prices); i++ {
		lastCapital, lastSold := _dp[0], _dp[1]
		_dp[0] = maximum(_dp[0], lastSold-prices[i])
		_dp[1] = maximum(_dp[1], lastCapital+prices[i]-fee)
	}

	return _dp[1]
}
