package dp

// https://leetcode.cn/problems/delete-operation-for-two-strings/

func MinDistance(word1, word2 string) int {
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	cols := len(word1)
	rows := len(word2)
	_dp := make([][]int, cols+1)
	for i := 0; i <= cols; i++ {
		_dp[i] = make([]int, rows+1)
		_dp[i][0] = i
	}
	for i := 0; i <= rows; i++ {
		_dp[0][i] = i
	}

	for i := 1; i <= cols; i++ {
		for j := 1; j <= rows; j++ {
			if word1[i-1] == word2[j-1] {
				_dp[i][j] = _dp[i-1][j-1]
			} else {
				_dp[i][j] = minimum(_dp[i-1][j], _dp[i][j-1]) + 1
			}
		}
	}

	return _dp[cols][rows]
}
