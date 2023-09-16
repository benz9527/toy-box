package b

// 前缀和

// 小明玩一个游戏。
// 系统发1+n张牌，每张牌上有一个整数。
// 第一张给小明，后n张按照发牌顺序排成连续的一行。
// 需要小明判断，后n张牌中，是否存在连续的若干张牌，其和可以整除小明手中牌上的数字。

func JudgeContinuousSumIsMultipleOfK(nums []int, m int) int {
	n := len(nums)
	preSum := make([]int, n+1)
	calcM := map[int]struct{}{}
	calcM[preSum[0]%m] = struct{}{}
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
		if _, ok := calcM[preSum[i]%m]; ok {
			return 1
		} else {
			calcM[preSum[i]%m] = struct{}{}
		}
	}
	return 0
}

// 给定一个小写字母组成的字符串 s，请找出字符串中两个不同位置的字符作为分割点，使得字符串分成三个连续子串且子串权重相等，
// 注意子串不包含分割点。
// 若能找到满足条件的两个分割点，请输出这两个分割点在字符串中的位置下标，若不能找到满足条件的分割点请返回0,0。
// 子串权重计算方式为:子串所有字符的ASCII码数值之和。
// 输入为一个字符串，字符串由a~z，26个小写字母组成，5 ≤ 字符串长度 ≤ 200。

func SameWeightSubContinuousStr(str string) [2]int {
	ans := [2]int{0, 0}
	prefixSum := make([]int, len(str)+1)
	for i, s := range str {
		prefixSum[i+1] = prefixSum[i] + int(s)

	}

	// i,j，i 不能从 0 开始，j 至少间隔 2
	for i := 1; i < len(str); i++ {
		weight1 := prefixSum[i] - prefixSum[0]
		for j := i + 2; j < len(str); j++ {
			weight2 := prefixSum[j] - prefixSum[i+1]
			if weight1 < weight2 {
				break
			}
			if weight1 == weight2 {
				weight3 := prefixSum[len(str)] - prefixSum[j+1]
				if weight2 == weight3 {
					return [2]int{i, j}
				}
			}
		}
	}

	return ans
}
