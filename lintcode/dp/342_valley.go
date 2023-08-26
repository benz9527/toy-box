package dp

// https://www.lintcode.com/problem/342/description
// 给你一个长度为 n的序列，在他的子序列中让你找一个山谷序列，山谷序列定义为：
// 1. 序列的长度为偶数。
// 2. 假设子序列的长度为 2 n。则前 n个数是严格递减的，后
//  n个数是严格递增的，并且第一段的最后一个元素和第二段的第一个元素相同，也是这段序列中的最小值。(不能使用同一个下标对应的元素作为中点)
// 3. 现在想让你找所有子序列中满足山谷序列规则的最长的长度为多少？
// 1<=len(num)<=1000
// 1<=num[i]<=10000
// 样例  1:
//	输入: num = [5,4,3,2,1,2,3,4,5]
//	输出: 8
//	样例解释:
//	最长山谷序列为[5,4,3,2,2,3,4,5]
//
//样例 2:
//	输入:  num = [1,2,3,4,5]
//	输出: 0
//	样例解释:
//	不存在山谷序列

func Valley(nums []int) int {
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	// _dp[i][0] 以 i 为终点，左侧递减的最大值
	// _dp[i][1] 以 i 为起点，右侧递增的最大值
	n := len(nums)
	if n <= 1 {
		return 0
	}

	_dp := make([][]int, n)
	for i := 0; i < n; i++ {
		_dp[i] = []int{1, 1} //
	}

	// left，这个子序列不是一个连续的操作，相当于可以删除中间的某些元素
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[i] < nums[j] {
				_dp[i][0] = maximum(_dp[i][0], _dp[j][0]+1)
			}
		}
	}
	// right
	for i := n - 2; i >= 0; i-- {
		for j := n - 1; j > i; j-- {
			if nums[i] < nums[j] {
				_dp[i][1] = maximum(_dp[i][1], _dp[j][1]+1)
			}
		}
	}
	res := 0
	for i := 0; i < n; i++ { // 先找左侧的最大递减子序列数
		for j := i + 1; j < n; j++ { // 再从右边找最大递增子序列数
			if nums[i] != nums[j] {
				continue
			}
			seq := minimum(_dp[i][0], _dp[j][1]) * 2
			if seq > res {
				res = seq
			}
		}
	}
	return res
}
