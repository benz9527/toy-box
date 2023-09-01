package backstrace

import isort "sort"

// https://leetcode.cn/problems/partition-to-k-equal-sum-subsets/
// 给定一个整数数组  nums 和一个正整数 k，找出是否有可能把这个数组分成 k 个非空子集，其总和都相等。

type descSlice []int

func (s descSlice) Less(i, j int) bool {
	return s[j] < s[i]
}
func (s descSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s descSlice) Len() int {
	return len(s)
}

func CanPartitionKSubsets(nums []int, k int) bool {
	sum := 0
	for i := 0; i < len(nums); i++ {
		sum += nums[i]
	}
	if sum%k > 0 {
		return false
	}
	subPivot := sum / k
	isort.Sort(descSlice(nums))
	if subPivot < nums[0] {
		return false
	}

	buckets := make([]int, k)
	var partition func(idx int) bool
	partition = func(idx int) bool {
		if idx == len(nums) {
			return true
		}
		num := nums[idx]
		for i := 0; i < k; i++ {
			if i > 0 && buckets[i] == buckets[i-1] {
				continue
			}
			// 标杆校准
			if num+buckets[i] <= subPivot {
				// 操作
				buckets[i] += num
				// 递归，下一个子问题
				if partition(idx + 1) {
					return true
				}
				// 逆操作
				buckets[i] -= num
			}
		}
		return false
	}
	return partition(0)
}
