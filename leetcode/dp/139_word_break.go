package dp

// https://leetcode.cn/problems/word-break/

// 对字符串的枚举内容作为数据来源（物品）
// 字符串的字串长度作为背包
// 字符串字串还需要继续切分为更小的字串进行判断
// 这个不是组合是排序，因为子字符串是有先后顺序的
// 比如 leetcode [leet,code]，前 4 个字符组成的串，会使得
// 字符串长度为 4 的背包填满，在长度为 8 的遍历时，分割到 code 的时候
// 其判断前 4 位的 leet 对应长度的背包被装满，code 才能装入
// 到长度为 8 的背包里面，因为本次遍历的长度就是 8

func WordBreak(s string, wordDict []string) bool {
	_dp := make([]bool, len(s)+1)
	_dp[0] = true
	set := map[string]struct{}{}
	for _, word := range wordDict {
		set[word] = struct{}{}
	}
	for i := 1; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			word := s[j:i]
			if _, ok := set[word]; ok && _dp[j] {
				_dp[i] = true
				break
			}
		}
	}

	return _dp[len(s)]
}

func WordBreakOptimize(s string, wordDict []string) bool {
	_dp := make([]bool, len(s)+1)
	_dp[0] = true
	set := map[string]struct{}{}
	mx, mi := 0, len(s)
	for _, word := range wordDict {
		set[word] = struct{}{}
		if len(word) > mx {
			mx = len(word)
		}
		if len(word) < mi {
			mi = len(word)
		}
	}
	for i := 1; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			if i-j > mx || i-j < mi {
				continue // 过滤无效的判断 4ms -> 0ms
			}
			word := s[j:i]
			if _, ok := set[word]; ok && _dp[j] {
				_dp[i] = true
				break
			}
		}
	}

	return _dp[len(s)]
}
