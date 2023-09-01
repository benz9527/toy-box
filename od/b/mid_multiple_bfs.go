package b

import (
	"math"
)

// 多源 bfs

// 在一个地图中(地图由n*n个区域组成），有部分区域被感染病菌。 感染区域每天都会把周围（上下左右）的4个区域感染。
// 请根据给定的地图计算，多少天以后，全部区域都会被感染。 如果初始地图上所有区域全部都被感染，或者没有被感染区域，返回-1

// 同类型题目
// 2XXX年，人类通过对火星的大气进行宜居改造分析，使得火星已在理论上具备人类宜居的条件；
// 由于技术原因，无法一次性将火星大气全部改造，只能通过局部处理形式；
// 假设将火星待改造的区域为row * column的网格，每个网格有3个值，宜居区、可改造区、死亡区，使用YES、NO、NA代替，YES表示该网格已经完成大气改造，NO表示该网格未进行改造，后期可进行改造，NA表示死亡区，不作为判断是否改造完的宜居，无法穿过；
// 初始化下，该区域可能存在多个宜居区，并目每个宜居区能同时在每个大阳日单位向上下左右四个方向的相邻格子进行扩散，自动将4个方向相邻的真空区改造成宜居区；
// 请计算这个待改造区域的网格中，可改造区是否能全部成宜居区，如果可以，则返回改造的大阳日天数，不可以则返回-1

func InfectionDays(area []int) int {
	n := int(math.Sqrt(float64(len(area))))
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = area[i*n+j]
		}
	}
	// queue [0] x [1] y
	queue := make([][2]int, 0, 16)
	lpop := func() [2]int {
		res := queue[0]
		queue = queue[1:]
		return res
	}
	rpush := func(pos ...[2]int) {
		queue = append(queue, pos...)
	}
	isEmpty := func() bool {
		return len(queue) == 0
	}
	ansDays := 0
	// 0 没感染 1 感染
	for rows := 0; rows < n; rows++ {
		for cols := 0; cols < n; cols++ {
			if area[rows*n+cols] == 1 {
				rpush([2]int{cols, rows})
			}
		}
	}
	// 第一天没有地方感染，或者全感染
	if isEmpty() || len(queue) == len(area) {
		return -1
	}

	nextAreas := [][2]int{
		{0, -1}, {0, 1}, {-1, 0}, {1, 0},
	}
	rest := len(area) - len(queue)

	for !isEmpty() && rest > 0 {
		nextDayAreas := make([][2]int, 0)
		for inflectedArea := lpop(); !isEmpty(); inflectedArea = lpop() {
			for _, nextArea := range nextAreas {
				x, y := inflectedArea[0]+nextArea[0], inflectedArea[1]+nextArea[1]
				if x >= 0 && x < n && y >= 0 && y < n && matrix[y][x] == 0 {
					nextDayAreas = append(nextDayAreas, nextArea)
					matrix[y][x] = 1
					rest--
				}
			}
		}
		ansDays++
		rpush(nextDayAreas...)
	}
	return ansDays
}

func InfectionDays2(area []int) int {
	n := int(math.Sqrt(float64(len(area))))
	// queue [0] x [1] y
	queue := make([][2]int, 0, 16)
	lpop := func() [2]int {
		res := queue[0]
		queue = queue[1:]
		return res
	}
	rpush := func(pos ...[2]int) {
		queue = append(queue, pos...)
	}
	isEmpty := func() bool {
		return len(queue) == 0
	}
	ansDays := 0
	// 0 没感染 1 感染
	for rows := 0; rows < n; rows++ {
		for cols := 0; cols < n; cols++ {
			if area[rows*n+cols] == 1 {
				rpush([2]int{cols, rows})
			}
		}
	}
	// 第一天没有地方感染，或者全感染
	if isEmpty() || len(queue) == len(area) {
		return -1
	}

	nextAreas := [][2]int{
		{0, -1}, {0, 1}, {-1, 0}, {1, 0},
	}
	rest := len(area) - len(queue)
	for !isEmpty() && rest > 0 {
		nextDayAreas := make([][2]int, 0)
		for inflectedArea := lpop(); !isEmpty(); inflectedArea = lpop() {
			for _, nextArea := range nextAreas {
				x, y := inflectedArea[0]+nextArea[0], inflectedArea[1]+nextArea[1]
				if x >= 0 && x < n && y >= 0 && y < n && area[y*n+x] == 0 {
					nextDayAreas = append(nextDayAreas, nextArea)
					area[y*n+x] = 1
					rest--
				}
			}
		}
		ansDays++
		rpush(nextDayAreas...)
	}
	return ansDays
}
