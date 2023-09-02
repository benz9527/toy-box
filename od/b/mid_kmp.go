package b

// KMP 字符串比较

// 在s中查找与t相匹配的子串，如果成功找到，则返回匹配的子串第一个字符在主串中的位置

func KMPIndexOf(src, tmpl string) int {
	next := KMPGetNextTable(tmpl)

	sIdx, tIdx := 0, 0
	for sIdx < len(src) && tIdx < len(tmpl) {
		if src[sIdx] == tmpl[tIdx] {
			// 同时往下走
			sIdx++
			tIdx++
		} else {
			if tIdx > 0 {
				// 只需要回退 tmpl 串的 tIdx 指针到 next[tIdx-1]位置，即最长相同前缀的结束位置后面一个位置
				tIdx = next[tIdx-1]
			} else {
				sIdx++
			}
		}
	}
	ans := -1
	if tIdx == len(tmpl) {
		// src 串中匹配 tmpl 的子串的首字符位置应该在 sIdx - t.length 位置，因为 sIdx
		// 指针最终会扫描到 src 串中匹配 tmpl 的子串的结束位置的后一个位置
		ans = sIdx - tIdx
	}
	return ans
}

func KMPGetNextTable(tmpl string) []int {
	next := make([]int, len(tmpl))
	j, k := 1, 0
	for j < len(tmpl) {
		if tmpl[j] == tmpl[k] {
			next[j] = k + 1
			j++
			k++
		} else {
			if k > 0 {
				k = next[k-1]
			} else {
				j++
			}
		}
	}
	return next
}

func GetMinLoopSubStr(n int, nums []int) []int {
	_getNextTable := func(n int, tmpl []int) []int {
		next := make([]int, n)
		j, k := 1, 0
		for j < n {
			if nums[j] == nums[k] {
				next[j] = k + 1
				j++
				k++
			} else {
				if k > 0 {
					k = next[k-1]
				} else {
					j++
				}
			}
		}
		return next
	}

	next := _getNextTable(n, nums)
	// 最长前后缀字串长度
	psLen := next[n-1]
	// 最小重复字串长度
	minLoopSubStrLen := 0
	if n%(n-psLen) == 0 {
		minLoopSubStrLen = n - psLen
	} else {
		minLoopSubStrLen = n
	}

	return nums[0:minLoopSubStrLen]
}
