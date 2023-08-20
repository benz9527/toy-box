package dp

// https://leetcode.cn/problems/coin-change-ii

// 组合不强调元素之间的顺序，排列强调元素之间的顺序。
// https://github.com/benz9527/toy-box/blob/396757f002cb02cab9c7ee1fe179f2045bfbff37/lintcode/dp/669_coins_dp.go

func change(amount int, coins []int) int {
	_dp := make([]int, amount+1)
	_dp[0] = 1

	for _, c := range coins {
		for j := c; j <= amount; j++ {
			_dp[j] += _dp[j-c]
		}
	}
	return _dp[amount]
}
