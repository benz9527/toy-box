package dp

// https://leetcode.cn/problems/longest-palindromic-subsequence/
// 给你一个字符串 s ，找出其中最长的回文子序列，并返回该序列的长度。
// 子序列定义为：不改变剩余字符顺序的情况下，删除某些字符或者不删除任何字符形成的一个序列。

func LongestPalindromeSubseq(s string) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	n := len(s)
	_dp := make([][]int, n)
	for i := 0; i < n; i++ {
		_dp[i] = make([]int, n)
		_dp[i][i] = 1
	}
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			if s[j] == s[i] {
				_dp[i][j] = _dp[i+1][j-1] + 2
			} else {
				_dp[i][j] = maximum(_dp[i][j-1], _dp[i+1][j])
			}
		}
	}
	return _dp[0][n-1]
}
