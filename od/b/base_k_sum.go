package b

import isort "sort"

type bigIntSlice []int64

func (s bigIntSlice) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s bigIntSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s bigIntSlice) Len() int {
	return len(s)
}

func KSum(nums []int64, k int, target int64) int {
	if k > len(nums) {
		return 0
	}
	isort.Sort(bigIntSlice(nums))
	return kSum(nums, k, 0, 0, target, 0)
}

func kSum(nums []int64, k, start, count int, target, sum int64) int {
	if k < 2 {
		return count
	}

	if k == 2 {
		return _2Sum(nums, start, count, target, sum)
	}

	for i := start; i < len(nums)-k; i++ {
		if nums[i] > 0 && sum+nums[i] > target {
			break
		}
		if i > start && nums[i] == nums[i-1] {
			continue
		}
		count = kSum(nums, k-1, i+1, count, target, sum+nums[i])
	}
	return count
}

func _2Sum(nums []int64, start, count int, target, sum int64) int {
	l, r := start, len(nums)-1
	for l < r {
		curSum := sum + nums[l] + nums[r]
		if target < curSum {
			r--
		} else if target > curSum {
			l++
		} else {
			count++
			for l+1 < r && nums[l] == nums[l+1] {
				l++
			}
			for r-1 > l && nums[r] == nums[r-1] {
				r--
			}
			l++
			r--
		}
	}
	return count
}
