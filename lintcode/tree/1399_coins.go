package tree

// 递归容易 Time Limit Exceeded

func TakeCoins(list []int, k int) int {
	mx := 0
	var decision func(int, []int, int)
	decision = func(start int, arr []int, k1 int) {
		if k1 == 0 {
			if mx < start {
				mx = start
			}
			return
		}
		h, t := arr[0], arr[len(arr)-1]
		if start+h > start+t {
			decision(start+h, arr[1:], k1-1)
		}
		decision(start+t, arr[0:len(arr)-1], k1-1)
	}

	decision(list[0], list[1:], k-1)
	decision(list[len(list)-1], list[0:len(list)-1], k-1)
	return mx
}

func TakeCoins2(list []int, k int) int {
	mx := 0
	n := len(list) + 1
	prefixSums := make([]int, n)
	for i := 1; i <= len(list); i++ {
		prefixSums[i] = prefixSums[i-1] + list[i-1]
	}

	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for i := 0; i <= k; i++ {
		l, r := i, k-i
		cur := prefixSums[n-1] - prefixSums[n-1-r] + prefixSums[l]
		mx = maximum(mx, cur)
	}

	return mx
}
