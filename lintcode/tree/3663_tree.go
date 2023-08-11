package tree

// https://www.lintcode.com/problem/3663/solution/77901?showListFe=true&page=1&problemTypeId=2&pageSize=50

// 无向树就是无向无环图
// 定理: 在一个连通无向无环图中，以任意节点出发所能到达的最远节点，一定是该图直径的端点之一。
//

// 第一次使用 bfs 找到层次最深的端点，再通过该最深的端点进行一次 bfs，获取树的层级数量 n，直径就是 n - 1
// The first time BFS is used to find the deepest endpoint, and then BFS is performed once
// through the deepest node to obtain the number of levels n of the tree, and the diameter is n-1

// Question:
// 给你一棵「 无向树 」，请你计算并返回它的「 直径 」，即这棵树上最长简单路径的 边数。
//
// 现用一个由所有 边 组成的数组 edges 来表示一棵无向树，其中 edges[i] = [u, v] 表示节点 u 和 v 之间的双向边。
//
// 树上的节点都使用 edges 中的数 {0, 1, ..., edges.length} 来作为标记，每个节点上的标记都是 唯一的。
//
// edges 会形成一棵无向树。
// 0<=edges.length<10^4
// edges[i][0]!=edges[i][1]
// 0<=edges[i][j]<=edges.length
//
// 示例1：
// 输入：edges = [[0,1],[0,2]]
// 输出：2
// 解释：edges 组成的无向树 1 - 0 - 2 是最长路径
//
// 示例2：
// 输入：edges = [[0,1],[1,2],[2,3],[1,4],[4,5]]
// 输出：4
// 解释：edges 组成的无向树 3 - 2 - 1 - 4 - 5 是最长路径

func UndirectedTreeDiameter(edges [][]int) int {
	degreeTable := make([]int, len(edges)+1)
	treeTable := make([][]int, len(edges)+1)

	for _, pair := range edges {
		degreeTable[pair[0]]++
		degreeTable[pair[1]]++
		if len(treeTable[pair[0]]) == 0 {
			treeTable[pair[0]] = []int{pair[1]}
		} else {
			treeTable[pair[0]] = append(treeTable[pair[0]], pair[1])
		}
		if len(treeTable[pair[1]]) == 0 {
			treeTable[pair[1]] = []int{pair[0]}
		} else {
			treeTable[pair[1]] = append(treeTable[pair[1]], pair[0])
		}
	}

	var (
		bfs            func(vertex int) [][]int
		getSubVertexes func(vertexes []int, traversedVertexes ...int) []int
		distinct       func(vertexes []int, parents ...int) []int
	)
	distinct = func(vertexes []int, parents ...int) []int {
		rest := make([]int, 0, len(vertexes))
		if len(parents) == 0 {
			return vertexes
		}
		for _, v := range vertexes {
			conflict := false
			for _, p := range parents {
				if v == p {
					conflict = true
				}
			}
			if !conflict {
				rest = append(rest, v)
			}
		}
		return rest
	}
	getSubVertexes = func(vertexes []int, traversedVertexes ...int) []int {
		subList := make([]int, 0, 16)
		for _, v := range vertexes {
			if len(treeTable[v]) > 0 {
				subList = append(subList, distinct(treeTable[v], traversedVertexes...)...)
			}
		}
		return subList
	}
	bfs = func(rootVertex int) [][]int {
		levels := make([][]int, 0, len(edges)+1)
		levels = append(levels, []int{rootVertex})
		traversedVertexes := make([]int, 0, len(edges)+1)
		traversedVertexes = append(traversedVertexes, rootVertex)
		for subVertexes := getSubVertexes([]int{rootVertex}, traversedVertexes...); len(subVertexes) > 0; subVertexes = getSubVertexes(subVertexes, traversedVertexes...) {
			levels = append(levels, append([]int{}, subVertexes...))
			traversedVertexes = append(traversedVertexes, subVertexes...)
		}
		return levels
	}

	levels := bfs(edges[0][0])
	result := bfs(levels[len(levels)-1][0])
	return len(result) - 1
}
