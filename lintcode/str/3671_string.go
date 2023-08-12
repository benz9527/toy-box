package str

import (
	isort "sort"
	"strings"
)

// https://www.lintcode.com/problem/3671/solution/77932?showListFe=true&page=1&problemTypeId=2&pageSize=50
// 现有一个近义词表 synonyms 和一个句子 text， 其中，synonyms 表中的每一个子项都是一些近义词对 ，你需要将句子 text 中每个在表中出现的单词用它的近义词来进行替换。

func ReplaceSynonyms(synonyms [][]string, text string) []string {

	if len(synonyms) == 0 {
		// 防御性代码
		return []string{text}
	}

	var (
		src      = strings.Split(text, " ")                         // 数据源处理
		table    = make([][]string, len(src))                       // 之后数据源要转换为“邻接表”的形式
		calc     = make([]int, len(src))                            // 记录每个单词的访问次数
		unionMap = make(map[int]map[string]struct{}, len(synonyms)) // 字符串并查集
		indexMap = make(map[string]int, len(synonyms))              // 辅助构造字符串并查集的数组，主要是能把字符串及其近义词统一映射为一个唯一的索引
		res      = make([]string, 0, 16)                            // 结果集
		obj      = make([]string, 0, len(src))                      // 构造结果的中间图
	)

	// 字符串并查集构造
	for _, syn := range synonyms {
		// 获取两个字符串是否有完成到唯一索引的映射
		i1, ok1 := indexMap[syn[0]]
		i2, ok2 := indexMap[syn[1]]
		if !ok1 && !ok2 {
			// 两个字符串都没有完成映射，直接构建新的唯一索引
			i1, i2 = len(indexMap), len(indexMap)
		} else if ok2 && !ok1 {
			// 第二个字符串完成映射，但是第一个没有，以第二个字符串的索引为准
			i1 = i2
		} else if ok1 && !ok2 {
			// 第一个字符串完成映射，但是第二个没有，以第一个字符串的索引为准
			i2 = i1
		} else {
			// 两个字符串都完成了索引映射，则以第一个索引（可以改进为以数量多的一方）为准
			ss := unionMap[i2]
			// 把第二个的索引映射的字符串集合并入到第一个的集合内
			for s := range ss {
				indexMap[s] = i1
				unionMap[i1][s] = struct{}{}
			}
			// 清理第二个索引映射的集合
			unionMap[i2] = make(map[string]struct{})
		}
		// 把两个字符串的索引映射到同一个索引上
		indexMap[syn[0]], indexMap[syn[1]] = i1, i1

		// 正常索引映射集合的元素添加流程
		if len(unionMap[i1]) == 0 {
			unionMap[i1] = make(map[string]struct{})
		}
		unionMap[i1][syn[1]] = struct{}{}
		unionMap[i1][syn[0]] = struct{}{}
	}

	// 数据源处理，映射为邻接表
	for i, s := range src {
		var ss map[string]struct{}
		// 获取数据源中当前字符串其近义词对应的并查集集合索引
		idx, ok := indexMap[s]
		if ok {
			// 如果存在，说明当前字符串需要进行近义词替换，获取近义词集合
			ss, ok = unionMap[idx]
		}
		if ok && len(ss) > 0 {
			// 从无序集合中还原为按字典序排列的字符串列表，并映射到当前字符串的索引上
			arr := make([]string, 0, len(ss))
			for k := range ss {
				arr = append(arr, k)
			}
			isort.Sort(isort.StringSlice(arr))
			table[i] = arr
		} else {
			// 没有近义词集合的，只需要构建当前字符串的单个元素数组
			table[i] = []string{s}
		}
		// 顺便构建第一个中间图结果
		obj = append(obj, table[i][0])
	}

	// 回溯需要用到的栈操作函数
	pop := func() int {
		if len(obj) == 0 {
			// 回溯完毕
			obj = []string{}
			return -1
		}
		i := len(obj) - 1
		obj = obj[0:i]
		// 只有出栈才算完成一次来回，只有这里才对访问加 1
		calc[i]++
		return i
	}
	push := func(s string) {
		obj = append(obj, s)
	}

	// 记录第一次的结果
	res = append(res, strings.Join(obj, " "))
	// 开始对中间图结果进行回溯，从栈顶开始
	for i := pop(); i >= 0; i = pop() {
		l := len(table[i])
		if l > 1 {
			c := calc[i] % l
			if calc[i] > 0 && c == 0 {
				// 近义词替换完成一轮，需要往下一处进行近义词替换
				continue
			}
			// 近义词替换
			push(table[i][c])
			// 补全剩余的字符串元素
			rest := len(src) - len(obj)
			// 这里使用索引的好处就是能直接根据长度就能获取原始数据进行替换
			for j := 0; j < rest; j++ {
				push(table[len(obj)][0])
			}
			// 存放结果
			res = append(res, strings.Join(obj, " "))
		}
	}

	return res
}
