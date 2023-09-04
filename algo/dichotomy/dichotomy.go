package dichotomy

// 1。能找到存在于数组中，必然是 middle index
// 2. 不能找到，不存在于数组中，需要返回其应该被插入的位置

func BinarySearchLastIndex(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := (l + r) >> 1
		if nums[mid] > target {
			r = mid - 1
		} else if nums[mid] < target {
			l = mid + 1
		} else {
			if mid == len(nums)-1 || nums[mid]-nums[mid+1] != 0 {
				// 这个方式就是找到最后一个等同于目标值的索引位置
				return mid
			}
			// 这样做是为了让相等范围进一步缩小，最后让 middle 值定位到最后的且等同于目标值的索引位置
			l = mid + 1
		}
	}
	// 没找到，返回插入的位置，以负数的形式
	// 获取到结果之后需要转换
	return -l - 1
}

func BinarySearchFirstIndex(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := (l + r) >> 1
		if nums[mid] > target {
			r = mid - 1
		} else if nums[mid] < target {
			l = mid + 1
		} else {
			if mid == 0 || nums[mid]-nums[mid-1] != 0 {
				// 这个方式就是找到第一个等同于目标值的索引位置
				return mid
			}
			// 这样做是为了让相等范围进一步缩小，最后让 middle 值定位到第一个且等同于目标值的索引位置
			r = mid - 1
		}
	}
	// 没找到，返回插入的位置，以负数的形式
	// 获取到结果之后需要转换
	return -l - 1
}
