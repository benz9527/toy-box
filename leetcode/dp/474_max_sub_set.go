package dp

// https://leetcode.cn/problems/ones-and-zeroes/
// 给你一个二进制字符串数组 strs 和两个整数 m 和 n 。
// 请你找出并返回 strs 的最大子集的长度，该子集中 最多 有 m 个 0 和 n 个 1 。
// 如果 x 的所有元素也是 y 的元素，集合 x 是集合 y 的 子集 。
// 输入：strs = ["10", "0001", "111001", "1", "0"], m = 5, n = 3
// 输出：4
// 解释：最多有 5 个 0 和 3 个 1 的最大子集是 {"10","0001","1","0"} ，因此答案是 4 。
// 其他满足题意但较小的子集包括 {"0001","1"} 和 {"10","1","0"} 。
// {"111001"} 不满足题意，因为它含 4 个 1 ，大于 n 的值 3 。

// 就是说着一个背包里面最多能装 m 个 0 和 n 个 1
// 通常 dp[i] 中装的都是一个值，装两种值得通常是 dp table 的形式
// 但是这个问题又应该是 01 背包问题，因为物品最多只能选一次

func FindMaxForm(strs []string, m int, n int) int {
	type thing struct {
		one  int
		zero int
	}
	things := make([]thing, 0, len(strs))
	for _, str := range strs {
		l := len(str)
		t := thing{one: 0, zero: 0}
		for _, ch := range str {
			if ch == '0' {
				t.zero++
			}
		}
		t.one = l - t.zero
		things = append(things, t)
	}

	// 初始化
	_dp := map[thing]int{}
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			_dp[thing{one: j, zero: m}] = 0
		}
	}

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for _, t := range things {
		for i := m; i >= t.zero; i-- {
			for j := n; j >= t.one; j-- {
				k := thing{one: j, zero: i}
				// 使用 map 的性能，在这里比二维数组要差很多，大概 40 倍的性能
				_dp[k] = maximum(_dp[k], _dp[thing{one: k.one - t.one, zero: k.zero - t.zero}]+1)
			}
		}
	}

	return _dp[thing{one: n, zero: m}]
}

func FindMaxFormOptimize(strs []string, m int, n int) int {
	type thing struct {
		one  int
		zero int
	}
	things := make([]thing, 0, len(strs))
	for _, str := range strs {
		l := len(str)
		t := thing{one: 0, zero: 0}
		for _, ch := range str {
			if ch == '0' {
				t.zero++
			}
		}
		t.one = l - t.zero
		things = append(things, t)
	}

	// 初始化
	_dp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		_dp[i] = make([]int, n+1)
	}

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for _, t := range things {
		for i := m; i >= t.zero; i-- {
			for j := n; j >= t.one; j-- {
				_dp[i][j] = maximum(_dp[i][j], _dp[i-t.zero][j-t.one]+1)
			}
		}
	}

	return _dp[m][n]
}
