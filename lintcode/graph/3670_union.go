package graph

import (
	isort "sort"
)

// 描述
// 在一个社交圈子中，有 n 个人，编号从 0 到 n - 1。现在有一份日志列表 logs，其中 logs[i] = [time, x, y] 表示 x 和 y 在 time 时刻成为朋友相互认识。
//
// 友谊是 相互且具有传递性 的。 也就是说：
//
// 相互性：当 a 和 b 成为朋友，那么 b 的朋友中也有 a
// 传递性：当 a 和 b 成为朋友，b 的朋友中有 c，那么 a 和 c 也会成为朋友相互认识
// 返回这个圈子中所有人之间都相互认识的最早时间。如果找不到最早时间，则返回 -1。
//
// 2≤n≤100
// 1≤logs.length≤10^4
// logs[1] !=logs[2]
// 在 logs 元素中，给定的时间 time 不会出现重复元素
//
// 示例1.
// 输入：logs = [[20220101,0,1],[20220109,3,4],[20220304,0,4],[20220305,0,3],[20220404,2,4]]
//      n = 6
// 输出： 20220404
// 解释：
//  time = 20220101，0 和 1 成为好友，关系为：[0,1], [2], [3], [4]
//  time = 20220109，3 和 4 成为好友，关系为：[0,1], [2], [3,4]
//  time = 20220304，0 和 4 成为好友，关系为：[0,1,3,4], [2]
//  time = 20220305，0 和 3 已经是好友
//  time = 20220404，2 和 4 成为好友，关系为：[0,1,2,3,4]，所有人都相互认识

func EarliestAcq(logs [][]int, n int) int {

	var (
		realSize = n
		sets     []map[int]struct{}               // 存放所有共同认识的人的索引，在同一个集合里相当于拥有相同的根或者祖先
		unionMap = make([]int, n)                 // 并查集，用于快速查找当前人所在的集合，或者所共同祖先的信息
		times    = make([]int, 0, len(logs))      // 排序后的时间列表
		logsMap  = make(map[int][]int, len(logs)) // 时间点对应的人
	)

	for i := 0; i < n; i++ {
		unionMap[i] = -1
	}
	for i := 0; i < len(logs); i++ {
		if logs[i][1] != logs[i][2] {
			times = append(times, logs[i][0])
			logsMap[logs[i][0]] = []int{logs[i][1], logs[i][2]}
		}
	}
	isort.Sort(isort.IntSlice(times))

	// 边的数量（关系数量）少于人头数，更新退出条件
	if len(times) < n {
		realSize = len(times)
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for i := 0; i < len(times); i++ {
		vals := logsMap[times[i]]
		bucket := 0
		if unionMap[vals[0]] == -1 && unionMap[vals[1]] == -1 {
			// 两者都没有共同的集合，也没有任意一个有被分配到过集合内，创建新的集合
			bucket = len(sets)
			sets = append(sets, make(map[int]struct{}, realSize))
		} else if unionMap[vals[0]] != unionMap[vals[1]] {
			// 两者有不同的集合信息，其中必有一个是有分配过集合
			if unionMap[vals[0]] == -1 || unionMap[vals[1]] == -1 {
				// 获取其中被分配过的集合，作为两者的共同集合
				bucket = max(unionMap[vals[0]], unionMap[vals[1]])
			} else if unionMap[vals[0]] >= 0 && unionMap[vals[1]] >= 0 {
				// 两者都分配过集合，选取其中最小索引的集合作为合并的对象
				bucket = min(unionMap[vals[0]], unionMap[vals[1]])
				// 索引大的集合被合并完数据后要清空，但是不要删除集合信息
				dropBucket := max(unionMap[vals[0]], unionMap[vals[1]])
				if dropBucket > 0 {
					// 被动合并数据
					for k, v := range sets[dropBucket] {
						sets[bucket][k] = v
					}
					// 清空，因为可能存在多个集合一旦移除索引，可能导致使用后续索引的数据出现越界
					sets[dropBucket] = map[int]struct{}{}
					// 全量更新指向的集合信息，因为其已经被合并和清空
					for j := 0; j < n; j++ {
						if unionMap[j] == dropBucket {
							unionMap[j] = bucket
						}
					}
				}
			}
		}
		// 合并信息
		unionMap[vals[0]], unionMap[vals[1]] = bucket, bucket
		// 集合更新，快速去重
		sets[bucket][vals[0]], sets[bucket][vals[1]] = struct{}{}, struct{}{}
		// 判断是否可以得出结果
		if len(sets[0]) >= realSize {
			return times[i]
		}
	}

	return -1
}

type logsSlice [][]int

func (x logsSlice) Len() int {
	return len(x)
}

func (x logsSlice) Less(i, j int) bool {
	return x[i][0] < x[j][0]
}

func (x logsSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func EarliestAcq2(logs [][]int, n int) int {

	var (
		parents = make([]int, n) // 并查集数据源
		setSize = make([]int, n) // 用于辅助判断数据的合并方向
		find    func(int) int
		union   func(x, y int)
	)
	// 原地排序
	isort.Sort(logsSlice(logs))
	// 初始化
	for i := 0; i < n; i++ {
		parents[i] = i
		setSize[i] = 1
	}

	find = func(x int) int {
		if parents[x] == x {
			return x
		}
		return find(parents[x])
	}
	union = func(x, y int) {
		px, py := find(x), find(y)
		// 小集合与大集合产生联系后，小集合忘大集合里合并
		if setSize[px] > setSize[py] {
			parents[py] = px
			setSize[px] += setSize[py]
		} else {
			parents[px] = py
			setSize[py] += setSize[px]
		}
	}

	// 退出条件更新
	if len(logs) < n {
		n = len(logs)
	}

	for i := 0; i < len(logs); i++ {
		ts, x, y := logs[i][0], logs[i][1], logs[i][2]
		px, py := find(x), find(y)
		// 有共同集合，继续执行
		if px == py {
			continue
		}
		// 合并没有共同集合的数据或者集合
		union(px, py)
		// 需要合并的次数应该逐渐减少，因为认识的人数在增加
		n--
		// n 个人之间，只需要认识不同人，应该是 n-1 次，就能所有人都认识了
		if n <= 1 {
			return ts
		}
	}

	return -1
}
