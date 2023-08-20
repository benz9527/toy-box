package dp

// https://leetcode.cn/problems/combination-sum-iv/

func CombinationSum4(nums []int, target int) int {
	_dp := make([]int, target+1)
	_dp[0] = 1
	// 组合问题，从小到大进行遍历（没有顺序要求，就是先遍历物品）
	// 但是这个问题中是有排序要求的，那么在循环中，要先遍历背包，再遍历物品
	// 先完成小背包的组合排列数，才能让大背包获得足够的组合数
	// 如果要把排列都列出来的话，只能使用回溯算法爆搜。
	for j := 0; j <= target; j++ {
		for _, n := range nums {
			if j < n {
				continue
			}
			_dp[j] += _dp[j-n]
		}
	}
	return _dp[target]
}
