package dp

// https://leetcode.cn/problems/coin-change
// https://github.com/benz9527/toy-box/blob/396757f002cb02cab9c7ee1fe179f2045bfbff37/lintcode/dp/669_coins_dp.go

func CoinChange(coins []int, amount int) int {
	_dp := make([]int, amount+1)
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	for i := 0; i <= amount; i++ {
		_dp[i] = -1
	}
	_dp[0] = 0
	for _, c := range coins {
		if c > amount {
			continue
		}
		for j := c; j <= amount; j++ {
			if j == c {
				_dp[j] = 1
			} else if j > c {
				if _dp[j-c] > 0 && _dp[j] < 0 {
					_dp[j] = _dp[j-c] + 1
				} else if _dp[j] > 0 && _dp[j-c] > 0 {
					_dp[j] = minimum(_dp[j], _dp[j-c]+1)
				}
			}
		}
	}
	return _dp[amount]
}

func CoinChange2(coins []int, amount int) int {
	_dp := make([]int, amount+1)
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	mx := 1<<32 - 1
	for i := 0; i <= amount; i++ {
		_dp[i] = 1<<32 - 1
	}
	_dp[0] = 0
	for _, c := range coins {
		if c > amount {
			continue
		}
		for j := c; j <= amount; j++ {
			if _dp[j-c] != mx {
				_dp[j] = minimum(_dp[j], _dp[j-c]+1)
			}
		}
	}
	if _dp[amount] == mx {
		return -1
	}
	return _dp[amount]
}
