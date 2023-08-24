package dp

// https://leetcode.cn/problems/maximum-length-of-repeated-subarray/

func MaxLengthOfRepeatedSubarray(nums1 []int, nums2 []int) int {
	rows := len(nums1)
	cols := len(nums2)
	_dp := make([][]int, cols+1)
	for i := 0; i <= cols; i++ {
		_dp[i] = make([]int, rows+1)
	}
	// 初始化所有都为 0
	res := 0
	for i := 1; i <= cols; i++ {
		for j := 1; j <= rows; j++ {
			if nums1[j-1] == nums2[i-1] {
				_dp[i][j] = _dp[i-1][j-1] + 1
			}
			if _dp[i][j] > res {
				res = _dp[i][j]
			}
		}
	}
	return res
}

func MaxLengthOfRepeatedSubarrayOptimize(nums1 []int, nums2 []int) int {
	rows := len(nums1)
	cols := len(nums2)
	_dp := make([]int, rows+1)
	// 初始化所有都为 0
	res := 0
	// 2倍的提升
	for i := 1; i <= cols; i++ {
		for j := rows; j > 0; j-- { // 避免重复覆盖
			if nums1[j-1] == nums2[i-1] {
				_dp[j] = _dp[j-1] + 1
			} else {
				_dp[j] = 0
			}
			if _dp[j] > res {
				res = _dp[j]
			}
		}
	}
	return res
}

func MaxLengthOfRepeatedSubarrayOptimize2(nums1 []int, nums2 []int) int {
	rows := len(nums1)
	cols := len(nums2)
	_dp := make([]int, rows+1)
	// 初始化所有都为 0
	res := 0
	for i := 1; i <= cols; i++ {
		last := make([]int, 0, rows+1)
		// 利用上一次的缓存去重，这里就是滚动数组的逻辑展开，
		// 也就是从后往前遍历所能够避免的重复或者覆盖问题
		last = append(last, _dp...)
		for j := 1; j <= rows; j++ {
			if nums1[j-1] == nums2[i-1] {
				_dp[j] = last[j-1] + 1
				if _dp[j] > res {
					res = _dp[j]
				}
			} else {
				_dp[j] = 0
			}
		}
	}
	return res
}
