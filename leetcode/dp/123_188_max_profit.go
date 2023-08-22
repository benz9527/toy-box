package dp

// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iii/

func MaxProfitIII(prices []int) int {

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
	_dp[0][1] = 0
	_dp[0][2] = -prices[0]
	_dp[0][3] = 0

	for i := 1; i < len(prices); i++ {
		_dp[i][0] = maximum(_dp[i-1][0], -prices[i])
		_dp[i][1] = maximum(_dp[i-1][1], _dp[i-1][0]+prices[i])
		_dp[i][2] = maximum(_dp[i-1][2], _dp[i][1]-prices[i])
		_dp[i][3] = maximum(_dp[i-1][3], _dp[i-1][2]+prices[i])
	}

	return _dp[len(prices)-1][3]
}

func MaxProfitIIIOptimize(prices []int) int {

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([][]int, 2)
	for i := 0; i < 2; i++ {
		_dp[i] = make([]int, 4)
	}

	// 第一次持有股票
	_dp[0][0] = -prices[0]
	// 第一次不持有股票（已抛售）
	_dp[0][1] = 0
	// 第二次持有股票
	_dp[0][2] = -prices[0] // 当天买卖两次，没有收益也没有成本，初始化应当假设它做了
	// 第二次不持有股票（已抛售）
	_dp[0][3] = 0

	for i := 1; i < len(prices); i++ {
		// 同只有一次买卖股票获取最大收益相同思路
		// 如果第一次买卖中已经获得最大收益了，第二次就当然还是这个最大收益
		// 这个思路是一次操作中尽可能模拟两次决策行为
		_dp[i%2][0] = maximum(_dp[(i-1)%2][0], -prices[i])
		_dp[i%2][1] = maximum(_dp[(i-1)%2][1], _dp[(i-1)%2][0]+prices[i])
		// 第二次以第一次抛售股票后的收益为基准进行股票买卖
		_dp[i%2][2] = maximum(_dp[(i-1)%2][2], _dp[i%2][1]-prices[i])
		_dp[i%2][3] = maximum(_dp[(i-1)%2][3], _dp[(i-1)%2][2]+prices[i])
	}

	return _dp[(len(prices)-1)%2][3]
}

func MaxProfitIIIOptimize2(prices []int) int {

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([]int, 4)

	_dp[0] = -prices[0]
	_dp[1] = 0
	_dp[2] = -prices[0]
	_dp[3] = 0

	for i := 1; i < len(prices); i++ {
		_dp[0] = maximum(_dp[0], -prices[i])
		_dp[1] = maximum(_dp[1], _dp[0]+prices[i])
		_dp[2] = maximum(_dp[2], _dp[1]-prices[i])
		_dp[3] = maximum(_dp[3], _dp[2]+prices[i])
	}

	return _dp[3]
}

// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iv/description/

// 给你一个整数数组 prices 和一个整数 k ，其中 prices[i] 是某支给定的股票在第 i 天的价格。
// 设计一个算法来计算你所能获取的最大利润。你最多可以完成 k 笔交易。也就是说，你最多可以买 k 次，卖 k 次。
// 注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。

func MaxProfitIV(k int, prices []int) int {

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([]int, 2*k)

	for i := 0; i < k; i++ {
		_dp[2*i] = -prices[0]
	}

	for i := 1; i < len(prices); i++ {
		_dp[0] = maximum(_dp[0], -prices[i])
		_dp[1] = maximum(_dp[1], _dp[0]+prices[i])
		for j := 1; j < k; j++ {
			_dp[2*j] = maximum(_dp[2*j], _dp[2*j-1]-prices[i])
			_dp[2*j+1] = maximum(_dp[2*j+1], _dp[2*j]+prices[i])
		}
	}

	return _dp[2*k-1]
}
