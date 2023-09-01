package b

// 01 背包
// 某组织举行会议，来了多个代表团同时到达，接待处只有一辆汽车，可以同时接待多个代表团，为了提高车辆利用率，请帮接待员计算可以坐满车的接待方案，输出方案数量。
// 约束:
// 一个团只能上一辆车，并且代表团人数 (代表团数量小于30，每个代表团人数小于30)小于汽车容量(汽车容量小于100)
// 需要将车辆坐满。

func FullCarForTravel(nums []int, n int) int {
	_dp := make([]int, n+1)
	_dp[0] = 1
	for i := 0; i < len(nums); i++ {
		for j := n; j >= nums[i]; j-- {
			_dp[j] += _dp[j-nums[i]]
		}
	}
	return _dp[n]
}
