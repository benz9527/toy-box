package graph

type DirectedGraphNode struct {
	Label     int
	Neighbors []*DirectedGraphNode
}

func NewDirectedGraphNode(x int) *DirectedGraphNode {
	return &DirectedGraphNode{
		Label:     x,
		Neighbors: []*DirectedGraphNode{},
	}
}

func (d *DirectedGraphNode) AddNeighbor(neighbors ...*DirectedGraphNode) *DirectedGraphNode {
	d.Neighbors = append(d.Neighbors, neighbors...)
	return d
}

// https://www.lintcode.com/problem/127/

func TopoSort(graph []*DirectedGraphNode) []*DirectedGraphNode {
	inDegreeMap := map[int]int{}
	outDegreeMap := map[int][]int{}
	for _, g := range graph {
		// 初始化入度图和出度图
		if _, ok := inDegreeMap[g.Label]; !ok {
			inDegreeMap[g.Label] = 0
		}
		if _, ok := outDegreeMap[g.Label]; !ok {
			outDegreeMap[g.Label] = []int{}
		}
		// 补全入度和出度关系
		for _, n := range g.Neighbors {
			if n.Label == g.Label {
				continue
			}
			if _, ok := inDegreeMap[n.Label]; !ok {
				inDegreeMap[n.Label] = 0
			}
			inDegreeMap[n.Label]++
			outDegreeMap[g.Label] = append(outDegreeMap[g.Label], n.Label)
		}
	}

	res := make([]*DirectedGraphNode, 0, len(inDegreeMap))
	for {
		inDegreeDescList := make([]int, 0, len(inDegreeMap))
		for k, inDegree := range inDegreeMap {
			if inDegree == 0 {
				// 以每一次遍历中入度为 0 的结点当做同一层
				res = append(res, &DirectedGraphNode{Label: k})
				inDegreeMap[k] = -1
				// 出度关系用于减少对应结点入度数量，相当于 BFS 中的寻找下一层
				out := outDegreeMap[k]
				if len(out) > 0 {
					inDegreeDescList = append(inDegreeDescList, out...)
				}
			}
		}
		if len(inDegreeDescList) > 0 {
			// 修正出度关系
			for _, nodeLabel := range inDegreeDescList {
				inDegreeMap[nodeLabel]--
			}
		} else {
			// 如果待减少入度的结点数量为 0，说明拓扑排序结束
			break
		}
	}
	// 如果最终参与排序的结点数量少于图中的结点数，说明图中存在环
	// 有向有环图，其入度数量是不可能减少的
	return res
}
