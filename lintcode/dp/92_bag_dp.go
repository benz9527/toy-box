package dp

func BackPack(m int, a []int) int {
	if m == 0 || len(a) == 0 {
		return 0
	}

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp := make([]int, m+1)
	for _, item := range a {
		for j := m; j >= item; j-- {
			_dp[j] = maximum(_dp[j-item]+item, _dp[j])
		}
	}
	return _dp[m]
}
