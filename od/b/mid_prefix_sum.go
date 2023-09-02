package b

// 前缀和

// 小明玩一个游戏。
// 系统发1+n张牌，每张牌上有一个整数。
// 第一张给小明，后n张按照发牌顺序排成连续的一行。
// 需要小明判断，后n张牌中，是否存在连续的若干张牌，其和可以整除小明手中牌上的数字。

func JudgeContinuousSumIsMultipleOfK(nums []int, m int) int {
	n := len(nums)
	preSum := make([]int, n)
	preSum[0] = nums[0]
	calcM := map[int]struct{}{}
	calcM[preSum[0]%m] = struct{}{}
	for i := 1; i < n; i++ {
		preSum[i] = preSum[i-1] + nums[i]
		if _, ok := calcM[preSum[i]%m]; ok {
			return 1
		} else {
			calcM[preSum[i]%m] = struct{}{}
		}
	}
	return 0
}
