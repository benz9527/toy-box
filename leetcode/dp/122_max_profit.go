package dp

// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-ii/

func MaxProfitII(prices []int) int {
	_dp := make([][]int, len(prices))
	for i := 0; i < len(prices); i++ {
		_dp[i] = make([]int, 2)
	}

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	// _d[i][0] 应当理解为前一天抛售股票所得利润+今天购入股票花费金额后的总值
	// _d[i][1] 应当理解为不断交易过程中所能获得的最大利润
	_dp[0][0] = -prices[0] // 一开始是要靠持有股票来进行，这第一天是否持股只是靠后面所得利润来判断
	_dp[0][1] = 0          // 一开始所得利润，最大就是不持有股票

	for i := 1; i < len(prices); i++ {
		//    7                1
		//
		//+------+         +------+
		//|      |         |      |
		//|   0  +-------->+   0  |
		//|      |         |      |
		//+-----<<o      XXX>------+
		//        oooooXX
		//           Xoo
		//        XXXX   oo
		//+------<X       o>>------+
		//|      |         |      |
		//|  -7  +-------->+  -1  |
		//|      |         |      |
		//+------+         +------+
		// _dp[i][0] 的来源
		// 1. 前一（或 N）天抛售股票所得收益 + 那一天购入股票的支出   (假设为 a)
		// 2. 当前这一天抛售股票所得收益 + 当前这一天购入股票的支出 (假设为 b)
		// 取两者的最大者
		// 如果 a >= b，说明当前价格购入股票会降低收益，故不进行购入，继续持有，也就不抛售
		// 如果 a < b，说明前一天抛售股票所得利润+当前价格购入股票仍然有更低的成本或者正向收益，支持购入和前一次完成抛售
		_dp[i][0] = maximum(_dp[i-1][0], _dp[i-1][1]-prices[i])
		// _dp[i][1] 的来源
		// 1. 前一（或 N）天的抛售股票所得的收益 (假设为 a)
		// 2. 前一天的持股成本 + 当前这一天股票的价格进行抛售后或的收益 (假设为 b)
		// 取两者的最大值
		// 如果 a >= b，收益低，不支持抛售股票
		// 如果 a < b, 收益有增加，支持抛售股票
		_dp[i][1] = maximum(_dp[i-1][1], prices[i]+_dp[i-1][0])
	}

	return _dp[len(prices)-1][1]
}

func MaxProfitIIOptimize(prices []int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	_dp := make([][]int, 2)
	for i := 0; i < 2; i++ {
		_dp[i] = make([]int, 2)
	}
	_dp[0][0] = -prices[0]
	_dp[0][1] = 0

	for i := 1; i < len(prices); i++ {
		_dp[i%2][0] = maximum(_dp[(i-1)%2][0], _dp[(i-1)%2][1]-prices[i])
		_dp[i%2][1] = maximum(_dp[(i-1)%2][1], prices[i]+_dp[(i-1)%2][0])
	}

	return _dp[(len(prices)-1)%2][1]
}

func MaxProfitIIOptimize2(prices []int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	_dp := make([]int, 2)
	_dp[0] = -prices[0]
	_dp[1] = 0

	for i := 1; i < len(prices); i++ {
		tmp := _dp[0]
		_dp[0] = maximum(_dp[0], _dp[1]-prices[i])
		_dp[1] = maximum(_dp[1], prices[i]+tmp)
	}

	return _dp[1]
}
