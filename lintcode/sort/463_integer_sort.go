package sort

import "github.com/benz9527/toy-box/algo/sort"

func SortIntegers(src []int) {
	l := len(src)
	for i := 0; i < l; i++ {
		next := i + 1
		if next >= l {
			break
		}
		if src[i] > src[next] {
			j := i
			for ; src[j] > src[next] && j > 0; j-- {
				temp := src[j]
				src[j] = src[next]
				src[next] = temp
				next--
			}
			if j == 0 && next == 1 && src[j] > src[next] {
				temp := src[j]
				src[j] = src[next]
				src[next] = temp
			}
		}
	}
}

func SortIntegersByAlgoBubble(src []int) {
	sort.BubbleSort(src)
}
