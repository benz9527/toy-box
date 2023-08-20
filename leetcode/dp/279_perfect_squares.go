package dp

// https://leetcode.cn/problems/perfect-squares/

// 给你一个整数 n ，返回 和为 n 的完全平方数的最少数量 。
// 完全平方数 是一个整数，其值等于另一个整数的平方；换句话说，其值等于一个整数自乘的积。
// 例如，1、4、9 和 16 都是完全平方数，而 3 和 11 不是。

func NumSquares(n int) int {
	_dp := make([]int, n+1)
	psList := make([]int, 0, 8)
	// 数据源要准确
	for i := 1; i*i <= n; i++ {
		psList = append(psList, i*i)
	}
	mx := 1<<32 - 1
	for i := 0; i <= n; i++ {
		_dp[i] = mx
	}
	// 初始化要准确
	_dp[0] = 0
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	for _, ps := range psList {
		for i := ps; i <= n; i++ {
			if _dp[i-ps] != mx {
				_dp[i] = minimum(_dp[i], _dp[i-ps]+1)
			}
		}
	}
	return _dp[n]
}
