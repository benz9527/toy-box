package dp

// https://leetcode.cn/problems/uncrossed-lines/
// 相对有序的子序列

func MaxUncrossedLines(nums1 []int, nums2 []int) int {
	cols := len(nums1)
	rows := len(nums2)
	_dp := make([][]int, cols+1)
	for i := 0; i <= cols; i++ {
		_dp[i] = make([]int, rows+1)
	}
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	for i := 1; i <= cols; i++ {
		for j := 1; j <= rows; j++ {
			if nums1[i-1] == nums2[j-1] {
				_dp[i][j] = _dp[i-1][j-1] + 1
			} else {
				_dp[i][j] = maximum(_dp[i-1][j], _dp[i][j-1])
			}
		}
	}
	return _dp[cols][rows]
}
