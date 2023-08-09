package recursion

func TowerOfHanoi(n int) []string {
	const (
		towerA = "A"
		towerB = "B"
		towerC = "C"
	)
	results := make([]string, 0, 128)
	return towerOfHanoi(n, towerA, towerC, towerB, results)
}

func towerOfHanoi(n int, from, to, dependOn string, results []string) []string {
	if n == 1 {
		results = append(results, move(1, from, to))
		return results
	}
	results = towerOfHanoi(n-1, from, dependOn, to, results)
	results = append(results, move(n, from, to))
	results = towerOfHanoi(n-1, dependOn, to, from, results)
	return results
}

func move(n int, from, to string) string {
	return "from " + from + " to " + to
}
