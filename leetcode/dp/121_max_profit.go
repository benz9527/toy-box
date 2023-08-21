package dp

// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock/

// Time limitation exceeded

func MaxProfit(prices []int) int {

	days := len(prices)
	_dp := make([]int, days)

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	mx := 0
	for i := 0; i < len(prices); i++ {
		for j := days - 1; j >= i; j-- {
			_dp[j] = maximum(_dp[j], prices[j]-prices[i])
			mx = maximum(mx, _dp[j])
		}
	}
	return mx
}

func MaxProfitOptimize(prices []int) int {
	days := len(prices)
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	buyDays := make([]int, 0, 16)  // 买点时间列表 (局部低点+潜在全局低点)
	saleDays := make([]int, 0, 16) // 买点时间列表 (局部高点+潜在全局高点)
	for day := 0; day < days; {
		curDay := day
		isTheLastDay := false
		// 处理股票上升的时间段
		for {
			nextDay := day + 1
			if nextDay < days && prices[nextDay] >= prices[day] {
				day++
			} else if nextDay == days {
				isTheLastDay = true
				break
			} else {
				break
			}
		}
		if curDay < day {
			buyDays = append(buyDays, curDay)
			saleDays = append(saleDays, day)
		} else {
			saleDays = append(saleDays, curDay)
		}
		if isTheLastDay {
			break
		}
		// 处理股票下降的时间段
		for {
			nextDay := day + 1
			if nextDay < days && prices[nextDay] < prices[day] {
				day++
			} else if nextDay == days {
				isTheLastDay = true
				break
			} else {
				break
			}
		}
		if isTheLastDay {
			break
		}
	}

	_dp := make([]int, len(buyDays))
	for j := 0; j < len(buyDays); j++ {
		for i := 0; i < len(saleDays); i++ {
			if buyDays[j] < saleDays[i] {
				_dp[j] = maximum(_dp[j], prices[saleDays[i]]-prices[buyDays[j]])
			}
		}
	}

	mx := 0
	for i := 0; i < len(buyDays); i++ {
		mx = maximum(mx, _dp[i])
	}
	return mx
}

func MaxProfitOptimize2(prices []int) int {
	days := len(prices)
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// 0: 持有股票（但是不卖出）获取的收益（成本）
	// 1: 未持有股票（曾经持有，但是当天卖出）获取的收益
	_dp := make([][]int, days)
	for i := 0; i < days; i++ {
		_dp[i] = make([]int, 2)
	}
	_dp[0][0] = -prices[0] // 假设这第一天开始持有该股票，其成本为当天的股票价格
	_dp[0][1] = 0          // 当天不持有股票，没有收益
	for j := 1; j < days; j++ {
		// 持有股票的收益
		// 1，继续持有原本股票成本
		// 2.（假设）第一次购入股票，只有股票成本比之前的持有成本低才算购入，否则是一直持有
		// 取两者之间的最大值
		_dp[j][0] = maximum(_dp[j-1][0], -prices[j])
		// 当前不持有股票的收益
		// 1. 当天卖出股票，收益为昨天的持有成本+当天的股价（该值较大，则是当天完成卖出行为）
		// 2. 之前就卖出了股票，收益为之前的收益 （该值较大，则是之前就完成了卖出行为）
		// 取两者之间的最大值
		_dp[j][1] = maximum(_dp[j-1][1], prices[j]+_dp[j-1][0])
	}
	return _dp[days-1][1]
}

func MaxProfitOptimize3(prices []int) int {
	days := len(prices)
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// 0: 持有股票（但是不卖出）获取的收益（成本）
	// 1: 未持有股票（曾经持有，但是当天卖出）获取的收益
	_dp := make([][]int, 2)
	for i := 0; i < 2; i++ {
		_dp[i] = make([]int, 2)
	}
	_dp[0][0] = -prices[0] // 假设这第一天开始持有该股票，其成本为当天的股票价格
	_dp[0][1] = 0          // 当天不持有股票，没有收益
	for j := 1; j < days; j++ {
		// 持有股票的收益
		// 1，继续持有原本股票成本
		// 2.（假设）第一次购入股票，只有股票成本比之前的持有成本低才算购入，否则是一直持有
		// 取两者之间的最大值
		_dp[j%2][0] = maximum(_dp[(j-1)%2][0], -prices[j])
		// 当前不持有股票的收益
		// 1. 当天卖出股票，收益为昨天的持有成本+当天的股价（该值较大，则是当天完成卖出行为）
		// 2. 之前就卖出了股票，收益为之前的收益 （该值较大，则是之前就完成了卖出行为）
		// 取两者之间的最大值
		_dp[j%2][1] = maximum(_dp[(j-1)%2][1], prices[j]+_dp[(j-1)%2][0])
	}
	return _dp[(days-1)%2][1]
}
