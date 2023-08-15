package arr

// https://leetcode.cn/problems/median-of-two-sorted-arrays/
// 寻找两个正序数组的中位数

func FindMedianSortedArrays(nums1 []int, nums2 []int) float64 {

	size := len(nums1) + len(nums2)
	maxIndex := size / 2
	halfArr := make([]int, 0, maxIndex+1)

	i, j := 0, 0
	for len(halfArr) <= maxIndex && i < len(nums1) && j < len(nums2) {
		n1, n2 := nums1[i], nums2[j]
		if n1 <= n2 {
			i++
			halfArr = append(halfArr, n1)
		} else {
			j++
			halfArr = append(halfArr, n2)
		}
	}
	if i >= len(nums1) && j < len(nums2) {
		res := maxIndex + 1 - len(halfArr)
		halfArr = append(halfArr, nums2[j:j+res]...)
	} else if i < len(nums1) && j >= len(nums2) {
		res := maxIndex + 1 - len(halfArr)
		halfArr = append(halfArr, nums1[i:i+res]...)
	}

	if size%2 == 1 {
		return float64(halfArr[len(halfArr)-1])
	}
	return float64(halfArr[len(halfArr)-2]+halfArr[len(halfArr)-1]) * 0.5
}
