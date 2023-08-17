package tree

// https://www.lintcode.com/problem/1399

// 递归容易 Time Limit Exceeded

func TakeCoins(list []int, k int) int {
	mx := 0
	var decision func(int, []int, int)
	// 这里分裂为两棵树，使用 DFS 的方式来穷举所有路径的权值之和，取其中的最大值
	decision = func(start int, arr []int, k1 int) {
		if k1 == 0 {
			if mx < start {
				mx = start
			}
			return
		}
		h, t := arr[0], arr[len(arr)-1]
		if start+h > start+t {
			// 以当前左边第一个元素为根，构建一棵树
			decision(start+h, arr[1:], k1-1)
		}
		// 以当前右边第一个元素为根，构建一棵树
		decision(start+t, arr[0:len(arr)-1], k1-1)
	}

	// 以左边第一个元素为根，构建一棵树
	decision(list[0], list[1:], k-1)
	// 以右边第一个元素为根，构建一棵树
	decision(list[len(list)-1], list[0:len(list)-1], k-1)
	return mx
}

// 使用前缀和进行查询优化，其时间复杂度是 O(1)，但是空间复杂度是 O(n)

func TakeCoins2(list []int, k int) int {
	mx := 0
	n := len(list) + 1
	prefixSums := make([]int, n)
	for i := 1; i <= len(list); i++ {
		prefixSums[i] = prefixSums[i-1] + list[i-1]
	}

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for i := 0; i <= k; i++ {
		l, r := i, k-i
		// 左右区间值的查询（枚举）
		// 右区间 + 左区间的值即为当前值
		cur := (prefixSums[n-1] - prefixSums[n-1-r]) + (prefixSums[0] + prefixSums[l])
		mx = maximum(mx, cur)
	}

	return mx
}
