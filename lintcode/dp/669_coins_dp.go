package dp

// https://www.lintcode.com/problem/669/description
// 给出不同面额的硬币以及一个总金额.
// 写一个方法来计算给出的总金额可以换取的最少的硬币数量.
// 如果已有硬币的任意组合均无法与总金额面额相等, 那么返回 -1.

// 可重复问题

func CoinChange(coins []int, amount int) int {
	optimizer := make([]int, amount+1)
	for i := 0; i <= amount; i++ {
		optimizer[i] = -1
	}
	optimizer[0] = 0
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	// 先遍历可选择的硬币，可以过滤掉不符合要求的，减少无用的步骤
	for _, coin := range coins {
		if coin > amount {
			continue
		}
		for subCoin := 1; subCoin <= amount; subCoin++ {
			if coin == subCoin {
				// 当硬币和子问题的值相等，直接最优（最少）的策略是 1
				optimizer[subCoin] = 1
			} else if coin < subCoin && optimizer[subCoin-coin] > 0 {
				// 只有当子问题的值大于当前的硬币值，才进行子问题的分裂
				// 直接从缓存中获取分裂后的子问题最优解
				if optimizer[subCoin] > 0 {
					// 当前子问题的最优解调整，从分裂子问题+1和当前记录的最优解中获取最小值
					optimizer[subCoin] = minimum(optimizer[subCoin-coin]+1, optimizer[subCoin])
				} else {
					// 当前子问题没有记录过最优解，直接从分裂子问题+1中获取最优解
					optimizer[subCoin] = optimizer[subCoin-coin] + 1
				}
			}
		}
	}
	// 直接查缓存
	return optimizer[amount]
}

func CoinChange2(coins []int, amount int) int {
	optimizer := make([]int, amount+1)
	for i := 0; i <= amount; i++ {
		optimizer[i] = -1
	}
	optimizer[0] = 0
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	// 先遍历子问题会让硬币的选择产生多余的无用执行次数，导致时间浪费
	for subCoin := 1; subCoin <= amount; subCoin++ {
		for _, coin := range coins {
			if coin > amount {
				continue
			}

			if coin == subCoin {
				optimizer[subCoin] = 1
			} else if coin < subCoin && optimizer[subCoin-coin] > 0 {
				if optimizer[subCoin] > 0 {
					optimizer[subCoin] = minimum(optimizer[subCoin-coin]+1, optimizer[subCoin])
				} else {
					optimizer[subCoin] = optimizer[subCoin-coin] + 1
				}
			}
		}
	}

	return optimizer[amount]
}
