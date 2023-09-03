package b

import isort "sort"

// 线段树

// 现在你要竞选一个县的县长。你去对每一个选民进行了调查。你已经知道每一个人要选的人是谁，以及要花多少钱才能让这个人选你。
// 现在你想要花最少的钱使得你当上县长。你当选的条件是你的票数比任何一个其它候选人的多（严格的多，不能和他们中最多的相等）。
// 请计算一下最少要花多少钱？
// 第一行有一个整数 n（1 ≤ n ≤ 10^5），表示这个县的选民数目。
//
// 接下来有 n 行，每一行有两个整数 ai 和 bi（0 ≤ ai ≤ 10^5；0 ≤ bi ≤ 10^4），
// 表示第 i 个选民选的是第 ai 号候选人，想要让他选择自己就要花 bi 的钱。你是 0 号候选人
// （所以，如果一个选民选你的话 ai 就是 0，这个时候 bi 也肯定是 0 ）。
//
// 把改选代价最大的选民留给其他候选人
// 扫描线，一次放弃一批，权值线段树用来存放被扫描后放弃的数据
// 利用权值线段树求解，这种时刻发生变化的集合的 前k小之和

type segmentNode struct {
	l, r, max int
	// max 只存这个 l-r区间内最大的数
}

type segmentTree struct {
	arr  []int
	tree []*segmentNode
}

func (s *segmentTree) build(k, l, r int) {
	// 初始化线段树节点, 即建立节点序号k和区间范围[l, r]的联系
	s.tree[k] = &segmentNode{l: l, r: r}
	if l == r {
		// 说明 k 节点是线段树的叶子节点
		// 线段树叶子节点的结果值就是arr[l]或arr[r]本身
		s.tree[k].max = s.arr[r]
		return
	}
	// 如果 l!=r, 则说明k节点不是线段树叶子节点，因此其必有左右子节点，左右子节点的分界位置是mid
	mid := (l + r) >> 1
	// 递归构建k节点的左子节点，序号为2 * k，对应区间范围是[l, mid]
	s.build(2*k, l, mid)
	// 递归构建k节点的右子节点，序号为2 * k + 1，对应区间范围是[mid+1, r]
	s.build(2*k+1, mid+1, r)
	//  k节点的结果值，取其左右子节点结果值的较大值
	s.tree[k].max = func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}(s.tree[2*k].max, s.tree[2*k+1].max)
}

func newSegmentTree(arr []int) *segmentTree {
	t := &segmentTree{
		// 执行查询区间最大值的原始数组
		arr: append([]int{}, arr...),
		// 线段树底层数据结构，以数组索引的关系映射为 tree，如果arr数组长度为n，则数组需要4n的长度
		// 因为是满二叉树，节点数 2 ^ (level-1) - 1
		// 补满最后一层 2 * n -1 + 2n 约等于 4n
		tree: make([]*segmentNode, len(arr)*4),
	}
	// 从根节点开始构建，线段树根节点序号k=1，对应的区间范围是[0, arr.length-1]
	t.build(1, 0, len(arr)-1)
	return t
}

// 权值线段树
type weightSegmentTree struct {
	weight []int // 权值树，统计的是某区间范围内元素的数量，这些数量累加起来就是对应元素的排名
	sum    []int // 和树，
}

func newWeightSegmentTree(size int) *weightSegmentTree {
	return &weightSegmentTree{
		weight: make([]int, size*4), // 初始化节点就是构建了这棵树
		sum:    make([]int, size*4),
	}
}

func (w *weightSegmentTree) addValue(k, l, r, val int) {
	if l == r {
		w.weight[k] += 1
		w.sum[k] += val
		return
	}

	// l != r
	mid := (l + r) >> 1
	kLeft, kRight := 2*k, 2*k+1
	if val <= mid {
		w.addValue(kLeft, l, mid, val)
	} else {
		w.addValue(kRight, mid+1, r, val)
	}
	w.weight[k] = w.weight[kLeft] + w.weight[kRight]
	w.sum[k] = w.sum[kLeft] + w.sum[kRight]
}

func (w *weightSegmentTree) query(k, l, r, rankWeight int) int {
	// 获取排序前 rankWeight 的和
	if l == r {
		return rankWeight * l
	}

	mid := (l + r) >> 1
	kLeft, kRight := 2*k, 2*k+1
	if w.weight[kLeft] < rankWeight {
		// 左子节点的元素数量 < rankWeight，那么说明前rank个数，还有部分在右子节点中
		return w.sum[kLeft] + w.query(kRight, mid+1, r, rankWeight-w.weight[kLeft])
	} else if w.weight[kLeft] > rankWeight {
		// 左子节点的元素数量 > rankWeight，那么说明前rank个数都在左子节点中，需要继续分解
		return w.query(kLeft, l, mid, rankWeight)
	}
	// 如果左子节点元素数量 == rankWeight，那么说明前 rankWeight 个数就是左子节点内的元素，
	// 此时要求前 rankWeight 小数之和，其实就是 w.sum[kLeft]
	return w.sum[kLeft]
}

func ElectionByMoney(peopleCount int, votes [][2]int) int {
	oneVoteMaxCost, ans := 0, 0
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	minimum := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	storeVoteMap := map[int][]int{}
	for i := 0; i < len(votes); i++ {
		store, cost := votes[i][0], votes[i][1]
		if cost == 0 || store == 1 {
			continue
		}
		oneVoteMaxCost = maximum(oneVoteMaxCost, cost)
		ans += cost
		_, ok := storeVoteMap[store]
		if !ok {
			storeVoteMap[store] = []int{}
		}
		storeVoteMap[store] = append(storeVoteMap[store], cost)
	}
	scanMatrix := make([][]int, 0, 16)
	for k, v := range storeVoteMap {
		isort.Sort(descSlice(v))
		storeVoteMap[k] = v
		// 横纵变换
		for i := 0; i < len(v); i++ {
			if len(scanMatrix) <= i {
				scanMatrix = append(scanMatrix, make([]int, 0, 16))
			}
			scanMatrix[i] = append(scanMatrix[i], v[i])
		}
	}
	// 一开始全部持有
	myVotes, myCosts := peopleCount, ans
	wst := newWeightSegmentTree(oneVoteMaxCost)
	for i := 0; i < len(scanMatrix); i++ {
		costs := scanMatrix[i]
		for j := 0; j < len(costs); j++ {
			// 扫描一次放弃这里所有的票
			wst.addValue(1, 1, oneVoteMaxCost, costs[j])
			// 由于放弃了这部分票数，要去除花费
			myCosts -= costs[j]
			// 由于放弃了这部分票数，要去除
			myVotes -= 1
		}
		extraVoteCost := 0
		// 此时其他店铺的选票数为 i+1，因此如果我们的选票数 myVotes <= i+1，则无法取胜
		if myVotes <= i+1 {
			// 手中总票数只有达到 i+2 张，才能战胜其他店铺，但是注意 i+2 不能超过总票数n,
			// 因此还需要正确额外的 min(i+2, n) - myVotes 张票
			extraVotes := minimum(i+2, peopleCount) - myVotes
			// 高效的求解一组数中的前x小数之和，可以基于权值线段树求解，相当于求解前 extraVotes 小之和
			extraVoteCost = wst.query(1, 1, oneVoteMaxCost, extraVotes)
		}
		// 每轮扫描线扫描后，计算其花费，保留最小花费作为题解
		ans = minimum(ans, myCosts+extraVoteCost)
	}
	return ans
}
