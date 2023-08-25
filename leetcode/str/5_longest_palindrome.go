package str

// https://leetcode.cn/problems/longest-palindromic-substring/

// 使用 Manacher's algo 的实现
// https://github.com/benz9527/toy-box/blob/7f4b3bdaf07695dacba5f658734d2980a2246d2c/lintcode/str/200_string.go

// 使用动态规划的实现

func LongestPalindrome(s string) string {
	n := len(s)
	ans := ""
	_dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		_dp[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := i; j < n; j++ {
			if s[i] == s[j] {
				if j-i <= 1 {
					_dp[i][j] = true
					tmp := s[i : j+1]
					if len(tmp) > len(ans) {
						ans = tmp
					}
				} else if _dp[i+1][j-1] {
					_dp[i][j] = true
					tmp := s[i : j+1]
					if len(tmp) > len(ans) {
						ans = tmp
					}
				}
			}
		}
	}
	return ans
}
