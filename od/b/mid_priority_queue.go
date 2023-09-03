package b

import (
	"container/heap"
	"errors"
	isort "sort"
)

// 最近，小明出了一些ACM编程题，决定在HDOJ举行一场公开赛。
// 假设题目的数量一共是n道，这些题目的难度被评级为一个不超过1000的非负整数，并且一场比赛至少需要一个题，而这场比赛的难度，就是所有题目的难度之和，同时，我们认为一场比赛与本场题目的顺序无关，而且题目也不会重复。
// 显而易见，很容易得到如下信息：
// 假设比赛只用1个题目，有n种方案；
// 假设比赛使用2个题目，有(n-1)*n/2种方案；
// 假设比赛使用3个题目，有(n-2)*(n-1)*n/6种方案；
// ............
// 假设比赛使用全部的n个题目，此时方案只有1种。
// 经过简单估算，小明发现总方案数几乎是一个天文数字！
// 为了简化问题，现在小明只想知道在所有的方案里面第m小的方案，它的比赛难度是多少呢？

// 商店里有N件唯一性商品，每件商品有一个价格，第 i 件商品的价格是 ai。
// 一个购买方案可以是从N件商品种选择任意件进行购买（至少一件），花费即价格之和。
// 现在你需要求出所有购买方案中花费前K小的方案，输出这些方案的花费。
// 当两个方案选择的商品集合至少有一件不同，视为不同方案，因此可能存在两个方案花费相同

var (
	globalPrices = []int{}
)

type combineModel struct {
	curSum, nextIdx int
}
type Item struct {
	Value *combineModel
	//Priority int
	index int
}
type queue []*Item

func (q queue) Len() int {
	return len(q)
}
func (q queue) Less(i, j int) bool {
	return q[i].Value.curSum+globalPrices[q[i].Value.nextIdx] < (q[j].Value.curSum + globalPrices[q[j].Value.nextIdx])
}
func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}
func (q *queue) Push(x interface{}) {
	n := len(*q)
	item := x.(*Item)
	item.index = n
	*q = append(*q, item)
}
func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*q = old[0 : n-1]
	return item
}

type PriorityQueue struct {
	data queue
}

func NewPriorityQueue() *PriorityQueue {
	pq := PriorityQueue{
		data: make(queue, 0),
	}
	heap.Init(&pq.data)
	return &pq
}

func (pq *PriorityQueue) Pop() (*Item, error) {
	if pq.data.Len() == 0 {
		return nil, errors.New("empty")
	}
	item := heap.Pop(&pq.data).(*Item)
	return item, nil
}

func (pq *PriorityQueue) Push(i *Item) {
	ni := &Item{
		Value: i.Value,
		//Priority: i.Priority,
	}
	heap.Push(&pq.data, ni)
}

func ShoppingPlans(n, k int, prices []int) []int {
	// n 商品个数
	// k 花费方案数
	// prices n 个商品的价格
	isort.Ints(prices)
	globalPrices = prices
	pq := NewPriorityQueue()
	zero := &combineModel{curSum: 0, nextIdx: 0}
	cur := &Item{Value: zero}
	ans := make([]int, 0, k)
	for i := 1; i <= k; i++ {
		ans = append(ans, cur.Value.curSum+prices[cur.Value.nextIdx])
		if cur.Value.nextIdx+1 < n {
			nc := &combineModel{curSum: cur.Value.curSum + prices[cur.Value.nextIdx], nextIdx: cur.Value.nextIdx + 1}
			nitem := &Item{Value: nc}
			pq.Push(nitem)
			cur.Value.nextIdx += 1
			pq.Push(cur)
		}
		cur, _ = pq.Pop()
	}
	return ans
}
