package dp

func LengthOfLIS(nums []int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	n := len(nums)
	// 存储到 i 以及之前子序列的最大升序长度
	_dp := make([]int, n)
	// 每个背包最短的最长子序列长度是 1
	for i := 0; i < n; i++ {
		_dp[i] = 1
	}
	// permutation problem
	res := _dp[0]            // 全局最长
	for i := 1; i < n; i++ { // knapsacks
		for j := 0; j < i; j++ { // sub knapsacks
			// 静态背包作为标杆去收集动态更小背包的结果
			// 相当于两次循环都在遍历背包
			// 这里进行小循环的逻辑是
			// nums[i] > nums[j], j 是从 0 -> i-1，nums[i] 就能从 nums[0] -> nums[i-1] 逐个比较大小
			// 这里的大小比较是间接传递的关系的
			// 比如 nums[i] > nums[j] 且 nums[j] > nums[j-1], j = i-1，那么 nums[i] > nums[j-1]，这种传递直到 nums[0] 为止。
			// 在这种 j 从 0 -> i-1 条件下 _dp[j] = _dp[j-1]+1，则 _dp[i] = _dp[j]+1。
			// 另一种条件下
			// 1. nums[i] > nums[j]，但是 nums[j] < nums[j-1]
			// 2. nums[i] < nums[j]，但是 nums[j] > nums[j-1]
			// 这两者会使得传递性中断，但是可以通过去掉一个不连续的结点
			// 让连续性继续保证，就是相当于取其中 _dp[j] 的最大值。
			// 当 i 越大，需要比较的子序列就越多，其中不连续递增的点就会越多，那么 _dp[i] 变换的最大值就可能会越多
			if nums[i] > nums[j] {
				// 最终就会是这个表达式（状态转移方程）
				_dp[i] = maximum(_dp[i], _dp[j]+1)
			}
		}
		if _dp[i] > res {
			res = _dp[i]
		}
	}
	return res
}
