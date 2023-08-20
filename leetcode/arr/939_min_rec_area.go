package arr

// https://leetcode.cn/problems/minimum-area-rectangle

// 低效穷举，对角线的端点，低效的原因在重复对角线运算

func MinAreaRect(points [][]int) int {
	mi := 1600000000
	type coordinate struct {
		y, x int
	}
	isDiagonal := func(c1, c2 coordinate) bool {
		return c1.x != c2.x && c1.y != c2.y
	}
	table := map[coordinate]struct{}{}
	restExists := func(c1, c2 coordinate) bool {
		_, ok1 := table[c1]
		_, ok2 := table[c2]
		return ok1 && ok2
	}
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}

	for i := 0; i < len(points); i++ {
		table[coordinate{y: points[i][1], x: points[i][0]}] = struct{}{}
	}

	for c1 := range table {
		for c2 := range table {
			if !isDiagonal(c1, c2) {
				continue
			}
			c3 := coordinate{y: c1.y, x: c2.x}
			c4 := coordinate{y: c2.y, x: c1.x}
			if !restExists(c3, c4) {
				continue
			}
			mi = minimum(mi, abs(c1.x-c2.x)*abs(c1.y-c2.y))
		}
	}

	if mi == 1600000000 {
		return 0
	}
	return mi
}
