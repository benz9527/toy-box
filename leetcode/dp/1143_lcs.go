package dp

// https://leetcode.cn/problems/longest-common-subsequence/

func LongestCommonSubsequence(text1, text2 string) int {
	cols := len(text1)
	rows := len(text2)
	_dp := make([][]int, cols+1)
	for i := 0; i <= cols; i++ {
		_dp[i] = make([]int, rows+1)
	}
	// text1 --> A, text2 --> B
	// _dp[i][j] = max(A[0:i-1], B[0:j-1])
	// 二维数组组成的区间中，最大的相同子序列的最大长度
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for i := 1; i <= cols; i++ {
		for j := 1; j <= rows; j++ {
			// 这里是有先后关系的，前一个 cols 的元素完成了最大相同子序列的规划
			// 后一个 cols 的元素只需要在前面的基础上进行就可以。
			// 同理 rows 的也是如此。此后有可能就是重复的计算
			// 因为这里是相对有序组成的子序列。
			// _dp[i][j] 的来源就有三个，同行左侧，同列上方以及左上角
			if text1[i-1] == text2[j-1] {
				_dp[i][j] = _dp[i-1][j-1] + 1
			} else {
				_dp[i][j] = maximum(_dp[i][j-1], _dp[i-1][j])
			}
		}
	}

	return _dp[cols][rows]
}
