package dp

// https://leetcode.cn/problems/longest-continuous-increasing-subsequence/

// 给定一个未经排序的整数数组，找到最长且 连续递增的子序列，并返回该序列的长度。
//
// 连续递增的子序列 可以由两个下标 l 和 r（l < r）确定，
// 如果对于每个 l <= i < r，都有 nums[i] < nums[i + 1] ，
// 那么子序列 [nums[l], nums[l + 1], ..., nums[r - 1],
// nums[r]] 就是连续递增子序列。

func FindLengthOfLCIS(nums []int) int {
	n := len(nums)
	_dp := make([]int, n)
	for i := 0; i < n; i++ {
		_dp[i] = 1
	}
	res := 1
	for i := 1; i < n; i++ {
		if nums[i-1] < nums[i] {
			_dp[i] += _dp[i-1]
		}
		if _dp[i] > res {
			res = _dp[i]
		}
	}
	return res
}
