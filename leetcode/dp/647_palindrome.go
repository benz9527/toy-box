package dp

// https://leetcode.cn/problems/palindromic-substrings/

func CountSubstring(s string) int {
	n := len(s)
	res := 0
	_dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		_dp[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := i; j < n; j++ {
			if s[i] == s[j] {
				if j-i <= 1 {
					res++
					_dp[i][j] = true
				} else if _dp[i+1][j-1] {
					res++
					_dp[i][j] = true
				}
			}
		}
	}
	return res
}
