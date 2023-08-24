package dp

// https://leetcode.cn/problems/distinct-subsequences/
// 困难

func NumDistinct(s, t string) int {
	cols := len(s)
	rows := len(t)
	_dp := make([][]int, cols+1)
	for i := 0; i <= cols; i++ {
		_dp[i] = make([]int, rows+1)
		_dp[i][0] = 1
	}
	for i := 1; i <= cols; i++ {
		for j := 1; j <= rows; j++ {
			if s[i-1] == t[j-1] {
				_dp[i][j] = _dp[i-1][j-1] + _dp[i-1][j]
			} else {
				_dp[i][j] = _dp[i-1][j]
			}
		}
	}

	return _dp[cols][rows]
}
