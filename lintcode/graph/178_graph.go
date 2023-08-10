package graph

func ValidTree(n int, edges [][]int) bool {
	if n == 0 {
		return false
	}
	if n == 1 {
		return true
	}
	if n > 1 {
		if len(edges) != n-1 || len(edges) > 0 && len(edges[0]) < 2 {
			return false
		}
	}

	unionMap := make(map[int]int, len(edges)*len(edges[0]))
	conflictEdges := make([][]int, 0, 8)

	find := func(v int) (int, bool) {
		ancestor, ok := unionMap[v]
		if !ok {
			return v, false
		}
		count := 0
		isLoop := false
		for {
			v, ok = unionMap[ancestor]
			if !ok {
				v = ancestor
				break
			}
			ancestor = v
			count++
			if count > n-1 {
				isLoop = true
				break
			}
		}
		return v, isLoop
	}

	union := func(v1, v2 int) bool {
		_, ok := unionMap[v2]
		if !ok {
			unionMap[v2] = v1
			return true
		}
		conflictEdges = append(conflictEdges, []int{v1, v2})
		return false
	}

	for _, pair := range edges {
		_ = union(pair[0], pair[1])
	}

	if len(conflictEdges) > 0 {
		for _, pair := range conflictEdges {
			anc, loop := find(pair[1])
			if loop {
				return false
			}
			op := unionMap[pair[1]]
			unionMap[pair[1]] = pair[0]
			anc2, loop2 := find(pair[1])
			if loop2 || anc == anc2 { // loop
				return false
			}
			unionMap[pair[1]] = op
		}
	} else {
		for _, pair := range edges {
			_, loop := find(pair[1])
			if loop {
				return false
			}
		}
	}

	return true
}

func ValidTree2(n int, edges [][]int) bool {
	if n == 0 {
		return false
	}
	if n == 1 {
		return true
	}
	if n > 1 {
		if len(edges) != n-1 || len(edges) > 0 && len(edges[0]) < 2 {
			return false
		}
	}

	unionMap := make(map[int]int, n)
	for i := 0; i < n; i++ {
		unionMap[i] = i
	}
	var find func(int) int
	find = func(v int) int {
		p := unionMap[v]
		if p == v {
			return v
		}
		return find(p)
	}

	union := func(v1, v2 int) bool {
		p1, p2 := find(v1), find(v2)
		if p1 == p2 {
			return false
		}
		unionMap[p1] = p2
		return true
	}

	for _, pair := range edges {
		if !union(pair[0], pair[1]) {
			return false
		}
	}

	return true
}

func ValidTree3(n int, edges [][]int) bool {
	if n == 0 {
		return false
	}
	if n == 1 {
		return true
	}
	if n > 1 {
		if len(edges) != n-1 || len(edges) > 0 && len(edges[0]) < 2 {
			return false
		}
	}

	unionMap := make([]int, n)
	for i := 0; i < n; i++ {
		unionMap[i] = i
	}
	var find func(int) int
	find = func(v int) int {
		p := unionMap[v]
		if p == v {
			return v
		}
		return find(p)
	}

	union := func(v1, v2 int) bool {
		p1, p2 := find(v1), find(v2)
		if p1 == p2 {
			return false
		}
		unionMap[p1] = p2
		return true
	}

	for _, pair := range edges {
		if !union(pair[0], pair[1]) {
			return false
		}
	}

	return true
}
