package dp

// https://leetcode.cn/problems/is-subsequence/
// 相对排序的子序列问题

func IsSubsequence(s, t string) bool {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	res := 0
	cols := len(s)
	rows := len(t)
	_dp := make([][]int, cols+1)
	for i := 0; i <= cols; i++ {
		_dp[i] = make([]int, rows+1)
	}
	for i := 1; i <= cols; i++ {
		for j := 1; j <= rows; j++ {
			if s[i-1] == t[j-1] {
				_dp[i][j] = _dp[i-1][j-1] + 1
			} else {
				_dp[i][j] = maximum(_dp[i-1][j], _dp[i][j-1])
			}
			if _dp[i][j] > res {
				res = _dp[i][j]
			}
		}
	}

	return res == cols
}
