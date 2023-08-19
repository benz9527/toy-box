package dp

// https://leetcode.cn/problems/unique-paths

// 一个机器人位于一个 m x n 网格的左上角 （起始点在下图中标记为 “Start” ）。
// 机器人每次只能向下或者向右移动一步。机器人试图达到网格的右下角（在下图中标记为 “Finish” ）

func UniquePaths(m, n int) int {
	y, x := m, n
	_dpTable := make([][]int, y)
	for i := 0; i < y; i++ {
		_dpTable[i] = make([]int, x) // 一开始全置零
		_dpTable[i][0] = 1           // 同一路径全为一
	}
	for j := 0; j < x; j++ {
		_dpTable[0][j] = 1 // 同一路径全为一
	}

	// 按增长顺序遍历就不会产生重复的子问题
	for i := 1; i < y; i++ {
		for j := 1; j < x; j++ {
			// 状态转移方程，就是靠移动的行为模式来确认的
			// 当前值只会来自其左边和上方的元素
			_dpTable[i][j] = _dpTable[i-1][j] + _dpTable[i][j-1]
		}
	}

	return _dpTable[y-1][x-1]
}
