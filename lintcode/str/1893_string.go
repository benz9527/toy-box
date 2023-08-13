package str

// https://www.lintcode.com/problem/1893/solution/77939?showListFe=true&page=2&problemTypeId=2&pageSize=50
// 描述
// 如果字符串的所有字符出现的次数相同，则认为该字符串是有效的。如果我们可以在字符串的某1个索引处删除1个字符，
// 并且其余字符出现的次数相同，那么它也是有效的。给定一个字符串s，判断它是否有效。如果是，返回YES，否则返回NO。
//
// 大部分情况下都是 NO，只有一部分符合条件的才是 YES
// 主要是靠出现次数进行统计。
//
// 如果字符出现次数种类大于两种，一定是 NO
// 如果字符出现次数种类等于一种，一定是 YES
// 如果字符出现次数种类是两种
//   - 两种次数之间的差值为1，只要任意一方的统计值为 1，一定是 YES
//   - 两种次数之间的差值大于 1，只能是次数种类值小的一方必须是次数种类值为 1 且统计值为 1，才是 YES

func IsValid(s string) string {
	m := map[int]int{}
	mx, mi := 0, 27
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return "NO"
		}
		idx := int(r - 'a')
		m[idx]++
		if mx < m[idx] {
			mx = m[idx]
		}
	}

	calc := map[int]int{}
	for _, v := range m {
		calc[v]++
	}

	if len(calc) == 1 {
		return "YES"
	}
	if len(calc) > 2 {
		return "NO"
	}

	for k := range calc {
		if k != mx {
			mi = k
		}
	}
	gap := mx - mi
	if gap == 1 && (calc[mx] == 1 || calc[mi] == 1) ||
		gap > 1 && calc[mi] == 1 && mi == 1 {
		return "YES"
	}

	return "NO"
}
