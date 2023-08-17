package dp

// 给出 n 个物品, 以及一个数组, nums[i] 代表第i个物品的大小, 保证大小均为正数,
// 正整数 target 表示背包的大小, 找到能填满背包的方案数。
// 每一个物品只能使用一次。
// 不可重复问题

func BackPackV(nums []int, target int) int {
	_dp := make([]int, target+1)
	// 初始化中，这里相当于 target 在分列为更小的子问题时
	// 如果子问题的值和输入数据中某个值相等，说明方案+1
	_dp[0] = 1 // 节省一个判断逻辑编写
	for i := 0; i < len(nums); i++ {
		// 遍历当前的输入数据，每次选择一个作为被拿掉的时候，
		// 判断剩余的数据能否合并当前的元素组合出一个方案
		// 假设 nums=[1,2,3,3,7], i=3 则 nums[i] = 3，
		// 给这个 3 一个标记，之后为 3<2>，那么 target = 7
		// 获取到 3<2> 之后，就剩 4 了，就要判断 4 这个子问题
		// 是否在之前的遍历中以及有过方案组合。
		// 假设 4 之前的组合方案有 1 + 3<1>，那么 7 在这次的遍历中
		// 就有 7 = 3<2> + 4 = 3<2> + 3<1> + 1 这一个组合方案。
		// 这种每个元素只能获取一次的操作，只能靠过去的迭代产生方案，不能关联
		// 未来的操作。基于过去产生现在的方案。
		for j := target; j >= nums[i]; j-- {
			_dp[j] += _dp[j-nums[i]]
		}
	}
	return _dp[target]
}

func BackPackV2(nums []int, target int) int {
	_dp := make([]int, target+1)
	for i := 0; i < len(nums); i++ {
		for j := target; j >= nums[i]; j-- {
			// _dp[0]=1 的逻辑展开
			if j-nums[i] == 0 {
				_dp[j] += 1
				continue
			}
			// _dp[j] += _dp[j-nums[i]] 的逻辑展开
			// 没有累积的方案可以选
			plans := _dp[j-nums[i]]
			if plans <= 0 {
				continue
			}
			_dp[j] += plans
		}
	}
	return _dp[target]
}
