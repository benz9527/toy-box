package arr

import isort "sort"

// https://leetcode.cn/problems/magnetic-force-between-two-balls/
// 在代号为 C-137 的地球上，Rick 发现如果他将两个球放在他新发明的篮子里，它们之间会形成特殊形式的磁力。Rick 有 n 个空的篮子，第 i 个篮子的位置在 position[i] ，Morty 想把 m 个球放到这些篮子里，使得任意两球间 最小磁力 最大。
// 已知两个球如果分别位于 x 和 y ，那么它们之间的磁力为 |x - y| 。
// 给你一个整数数组 position 和一个整数 m ，请你返回最大化的最小磁力。
//
// huawei od

func MaxDistanceBetweenBalls(position []int, m int) int {
	isort.Ints(position)
	maxDist := position[len(position)-1] - position[0]
	ans, minDist := 0, 0
	isOkDist := func(potentialDist int) bool {
		curPosition := position[0]
		count := 1
		for i := 1; i < len(position); i++ {
			if position[i]-curPosition >= potentialDist {
				count++
				curPosition = position[i]
			}
		}
		return count >= m
	}
	// 二分法退出条件：左侧边界大于右侧边界
	for minDist <= maxDist {
		midDist := (minDist + maxDist) >> 1
		if isOkDist(midDist) {
			ans = midDist
			minDist = midDist + 1
		} else {
			maxDist = midDist - 1
		}
	}

	return ans
}
