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

// 打家劫舍

// 小明和朋友玩跳格子游戏，有 n 个连续格子，每个格子有不同的分数，小朋友可以选择以任意格子起跳，但是不能跳连续的格子，也不能回头跳；
// 给定一个代表每个格子得分的非负整数数组，计算能够得到的最高分数。

func JumpGrids(grids []int) int {
	_dp := make([]int, len(grids))
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	_dp[0] = grids[0]
	if len(grids) > 1 {
		_dp[1] = maximum(_dp[0], grids[1])
	}
	for i := 2; i < len(grids); i++ {
		for j := len(grids) - 1; j >= i; j-- {
			_dp[j] = maximum(_dp[j-1], _dp[j-2]+grids[j])
		}
	}
	return _dp[len(grids)-1]
}

func JumpGridsII(grids []int) int {
	if len(grids) <= 0 {
		return 0
	}
	if len(grids) == 1 {
		return grids[0]
	}
	grids1 := grids[0 : len(grids)-1]
	grids2 := grids[1:]
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	_jumpGrids := func(grids []int) int {
		if len(grids) <= 0 {
			return 0
		}
		_dp := make([]int, len(grids))
		_dp[0] = grids[0]
		if len(grids) > 1 {
			_dp[1] = maximum(_dp[0], grids[1])
		}
		for i := 2; i < len(grids); i++ {
			for j := len(grids) - 1; j >= i; j-- {
				_dp[j] = maximum(_dp[j-1], _dp[j-2]+grids[j])
			}
		}
		return _dp[len(grids)-1]
	}
	return maximum(_jumpGrids(grids1), _jumpGrids(grids2))
}

// 平分问题的动态规划解
// MELON有一堆精美的雨花石（数量为n，重量各异），准备送给S和W。
// MELON希望送给俩人的雨花石重量一致，请你设计一个程序，帮MELON确认是否能将雨花石平均分配。
// 第1行输入为雨花石个数：n， 0 < n < 31。
// 第2行输入为空格分割的各雨花石重量：m[0] m[1] ….. m[n - 1]， 0 < m[k] < 1001。
// 不需要考虑异常输入的情况。

func Melon(stones []int) int {
	half := 0
	for i := 0; i < len(stones); i++ {
		half += stones[i]
	}
	if half&0x1 == 0x1 {
		return -1
	}
	half = half >> 1
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	_dp := make([][]int, len(stones)+1)
	for i := 0; i <= len(stones); i++ {
		_dp[i] = make([]int, half+1)
	}
	for i := 0; i <= half; i++ {
		_dp[0][i] = len(stones)
	}
	_dp[0][0] = 0
	for i := 1; i <= len(stones); i++ {
		stone := stones[i-1]
		for j := 1; j <= half; j++ {
			if stone < j {
				_dp[i][j] = minimum(_dp[i-1][j], _dp[i-1][j-stone]+1)
			} else if stone == j {
				_dp[i][j] = 1
			} else {
				_dp[i][j] = _dp[i-1][j]
			}
		}
	}
	ans := _dp[len(stones)][half]
	if _dp[len(stones)][half] == len(stones) {
		ans = -1
	}
	return ans
}

func Melon2(stones []int) int {
	half := 0
	for i := 0; i < len(stones); i++ {
		half += stones[i]
	}
	if half&0x1 == 0x1 {
		return -1
	}
	half = half >> 1
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	_dp := make([]int, half+1)
	for i := 0; i <= half; i++ {
		_dp[i] = len(stones)
	}
	_dp[0] = 0
	for i := 1; i <= len(stones); i++ {
		stone := stones[i-1]
		for j := 1; j <= half; j++ {
			if stone < j {
				_dp[j] = minimum(_dp[j], _dp[j-stone]+1)
			} else if stone == j {
				_dp[j] = 1
			} else {
				_dp[j] = _dp[j]
			}
		}
	}
	ans := _dp[half]
	if _dp[half] == len(stones) {
		ans = -1
	}
	return ans
}
