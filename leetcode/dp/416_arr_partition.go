package dp

// https://leetcode.cn/problems/partition-equal-subset-sum/

func CanPartition(nums []int) bool {
	sum := 0
	for i := 0; i < len(nums); i++ {
		sum += nums[i]
	}
	if sum&0x1 == 1 {
		return false
	}
	half := sum / 2
	// 这里背包是做装满次数计数，只要大于等于 2 就说明能够对半划分
	// 大于的原因是相同的数字可能重复出现导致重复计数
	_dp := make([]int, half+1)
	for _, n := range nums {
		for j := half; j >= n; j-- {
			if j-n == 0 || _dp[j-n] > 0 {
				_dp[j] += 1 // 防止累加溢出
			}
		}
	}
	return _dp[half] >= 2
}

func CanPartitionOptimize(nums []int) bool {
	sum := 0
	for i := 0; i < len(nums); i++ {
		sum += nums[i]
	}
	if sum&0x1 == 1 {
		return false
	}
	half := sum / 2

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// 01 knapsack
	_dp := make([]int, half+1)
	for _, n := range nums {
		for j := half; j >= n; j-- {
			_dp[j] = maximum(_dp[j], _dp[j-n]+n)
		}
	}
	return _dp[half] == half
}
