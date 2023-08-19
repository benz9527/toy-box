package dp

func LastStoneWeightII(stones []int) int {
	sum := 0
	for i := 0; i < len(stones); i++ {
		sum += stones[i]
	}
	half := sum / 2

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([]int, half+1)
	for _, stone := range stones {
		for i := half; i >= stone; i-- {
			// 这里最大的背包不一定能够装满，
			// 不满的情况下进行碰撞之后
			// 剩下的就是石头的最小值
			_dp[i] = maximum(_dp[i], _dp[i-stone]+stone)
		}
	}

	return sum - 2*_dp[half]
}
