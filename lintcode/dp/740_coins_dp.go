package dp

// https://www.lintcode.com/problem/740

// 给出不同面值的硬币以及总金额. 试写一函数来计算构成该总额的组合数量. 你可以假设每一种硬币你都有无限个.
// 0 <= amount <= 5000
// 1 <= coin <= 5000

func Change(amount int, coins []int) int {
	_dp := make([]int, amount+1)

	for i := 0; i < len(coins); i++ {
		if coins[i] > amount {
			continue
		}
		// 假设当前的硬币面值是 1，amount 为 5
		// 这里遍历的起点就是 _dp[1] += 1
		// 之后逐步递增，直到 amount
		// _dp[2] += _dp[1] : 1
		// _dp[3] += _dp[2] : 1
		// _dp[4] += _dp[3] : 1
		// _dp[5] += _dp[4] : 1
		// 上面的过程表明当前的组合策略都是只有 1 为元素的组合
		// 同理当面值为 2 时，遍历的起点就是 _dp[2] += 1，并不会再包括对 1 的元素更新
		// 这样就不会产生重复的计算
		for subQ := coins[i]; subQ <= amount; subQ++ {
			if subQ-coins[i] == 0 {
				// _dp[0] = 1 的逻辑展开
				_dp[subQ] += 1
			} else {
				// 计算当前子问题 subQ 获取到面值 x 之后，其差值 subQ-x 对应
				// 的问题上是否能增加组合数量。
				// subQ-x 的（子）问题的解就是档期可以增加的组合数
				// 因为 _dp[x] 所代表的它之前有过的组合数，现在 x 不变化，要 subQ-x 进行
				// 拆分演变，而 subQ-x 的拆分演变数量就是它到目前为止积累下来的组合数
				_dp[subQ] += _dp[subQ-coins[i]]
			}
		}
	}

	return _dp[amount]
}
