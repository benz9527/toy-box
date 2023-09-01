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

// 有一名科学家想要从一台古董电脑中拷贝文件到自己的电脑中加以研究。
// 但此电脑除了有一个3.5寸软盘驱动器以外，没有任何手段可以将文件持贝出来，而且只有一张软盘可以使用。
// 因此这一张软盘是唯一可以用来拷贝文件的载体。
// 科学家想要尽可能多地将计算机中的信息拷贝到软盘中，做到软盘中文件内容总大小最大。
// 已知该软盘容量为1474560字节。文件占用的软盘空间都是按块分配的，每个块大小为512个字节。
// 一个块只能被一个文件使用。拷贝到软盘中的文件必须是完整的，且不能采取任何压缩技术。

func MaxCopyFileSize(files []int, n int) int {
	const maxSize = 1474560
	// weight --> block
	// value --> byte
	bagSize := maxSize / 512
	_dp := make([]int, bagSize+1)
	fileBlocks := make([]int, n)
	for i := 0; i < n; i++ {
		fileBlocks[i] = files[i] / 512
		if files[i]%512 > 0 {
			fileBlocks[i]++
		}
	}

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	_dp[0] = 0
	for i := 0; i < n; i++ {
		for j := bagSize; j >= fileBlocks[i]; j-- {
			_dp[j] = maximum(_dp[j], _dp[j-fileBlocks[i]]+files[i])
		}
	}
	return _dp[bagSize]
}
