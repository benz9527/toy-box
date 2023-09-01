package b

import isort "sort"

type descSlice []int

func (s descSlice) Less(i, j int) bool {
	return s[j] < s[i]
}
func (s descSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s descSlice) Len() int {
	return len(s)
}

type descSlice2 [][]int

func (s descSlice2) Less(i, j int) bool {
	l1, l2 := s[i], s[j]
	for i := 0; i < len(l1); i++ {
		if l1[i] > l2[i] {
			return false
		}
		if l1[i] < l2[i] {
			return true
		}
		// ==
	}
	return false
}
func (s descSlice2) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s descSlice2) Len() int {
	return len(s)
}

func sumAll(nums ...int) int {
	sum := 0
	if len(nums) < 0 {
		return sum
	}
	for _, n := range nums {
		sum += n
	}
	return sum
}

func SecretiveElevator(nums []int, target, n int) []int {
	totalSum := sumAll(nums...)
	// 最大可以上升的次数
	totalUpTimes := n>>1 + n&0x1
	// up - (total - up) <= target
	// up <= (total + target) / 2
	// 最大上升范围
	upSumCeil := (totalSum + target) >> 1
	// 最优上升范围，不断更新，但不会超过限制
	upSumFloor := 0
	// 可行解
	ansPaths := [][]bool{}
	// 降序
	isort.Sort(descSlice(nums))
	// 深度遍历子树，找解
	var dfs func(nums []int, idx int, pathHelper []bool, sum, count int)
	dfs = func(nums []int, idx int, pathHelper []bool, sum, count int) {
		if count > totalUpTimes {
			return
		}
		if count == totalUpTimes {
			// 非最优解
			if sum < upSumFloor {
				return
			}
			// 发现更优解
			if sum > upSumFloor {
				upSumFloor = sum
				ansPaths = [][]bool{}
			}
			ansPaths = append(ansPaths, append([]bool{}, pathHelper...))
			return
		}
		for i := idx; i < len(nums); i++ {
			item := nums[i]
			if sum+item > upSumCeil {
				continue
			}
			pathHelper[i] = true
			dfs(nums, i+1, pathHelper, sum+item, count+1)
			pathHelper[i] = false
		}
	}
	mixup := func(nums []int, pathHelper []bool) []int {
		up, down, ans := []int{}, []int{}, []int{}
		for i := 0; i < len(nums); i++ {
			if pathHelper[i] {
				up = append(up, nums[i])
			} else {
				down = append(down, nums[i])
			}
		}
		for i := 0; i < len(nums)/2; i++ {
			ans = append(ans, up[0])
			up = up[1:]
			ans = append(ans, down[0])
			down = down[1:]
		}
		if len(up) > 0 {
			ans = append(ans, up[0])
		}
		return ans
	}
	dfs(nums, 0, make([]bool, len(nums)), 0, 0)
	if len(ansPaths) < 0 {
		return []int{}
	}
	ansList := [][]int{}
	for i := 0; i < len(ansPaths); i++ {
		ansList = append(ansList, mixup(nums, ansPaths[i]))
	}
	isort.Sort(descSlice2(ansList))
	return ansList[0]
}
