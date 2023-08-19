package dp

// https://leetcode.cn/problems/unique-paths-ii

func UniquePathsWithObstacles(obstacleGrid [][]int) int {
	y, x := len(obstacleGrid), len(obstacleGrid[0])
	_dpTable := make([][]int, y)
	for i := 0; i < y; i++ {
		_dpTable[i] = make([]int, x) // 一开始全置零
	}
	for i := 0; i < y; i++ {
		if i > 0 && obstacleGrid[i-1][0] == 1 {
			break
		}
		_dpTable[i][0] = 1
	}
	for j := 0; j < x; j++ {
		if j > 0 && obstacleGrid[0][j-1] == 1 {
			break
		}
		_dpTable[0][j] = 1
	}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if obstacleGrid[i][j] == 1 {
				_dpTable[i][j] = 0
			}
		}
	}

	getValueIfObstacle := func(i, j int) int {
		if obstacleGrid[i][j] == 1 {
			return 0
		}
		return _dpTable[i][j]
	}

	// 按增长顺序遍历就不会产生重复的子问题
	for i := 1; i < y; i++ {
		for j := 1; j < x; j++ {
			if obstacleGrid[i][j] == 1 {
				continue
			}
			_dpTable[i][j] = getValueIfObstacle(i-1, j) + getValueIfObstacle(i, j-1)
		}
	}

	return _dpTable[y-1][x-1]
}
