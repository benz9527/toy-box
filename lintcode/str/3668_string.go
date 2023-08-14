package str

import (
	isort "sort"
	"strings"
)

// https://www.lintcode.com/problem/3668/description?showListFe=true&page=1&problemTypeId=2&pageSize=50

// 给定一个字符串 s，按以下规则组成一个新的单词：
//
// 若字母在花括号之外，则说明该字母在新单词的对应位置是一个必选字符
// 若字母在花括号之内，则说明该字母在新单词的对应位置是一个选项字符
// 花括号之内可能会存在多个字母，对于花括号内的所有选项字符，必须且只能选取其中一个
// 按照以上规则，返回所有可以组成的新单词，并按照字典顺序进行排列。
//
// 1≤s.length≤100
// s 仅由括号 {}，逗号 , 和小写英文字母组成，不含有空格
// s 保证是一个有效的输入
// s 不含有嵌套括号
// 对于 s 中的每一组括号，在同一个括号内不会存在重复的字母

func BracketExpansion(s string) []string {

	var (
		res   []string
		table [][]string
		obj   []string
		calc  []int
	)

	// 使用邻接表构建树（数据源）
	for i := 0; i < len(s); i++ {
		if s[i] != '{' && s[i] != '}' && s[i] != ',' {
			table = append(table, []string{string(s[i])})
		} else if s[i] == '{' {
			collector := make([]string, 0, 4)
			for i++; s[i] != '}'; i++ {
				if s[i] == ',' {
					continue
				}
				collector = append(collector, string(s[i]))
			}
			isort.Sort(isort.StringSlice(collector))
			table = append(table, collector)
		}
	}

	// 对应层级的节点访问次数记录，用于辅助
	calc = make([]int, len(table))
	// 回溯时要使用的栈操作函数
	pop := func() int {
		if len(obj) == 0 {
			obj = []string{}
			return -1
		}
		idx := len(obj) - 1
		obj = obj[0:idx]
		// 一个来回算是完成一次回溯访问
		calc[idx]++
		return idx
	}
	push := func(str string) {
		obj = append(obj, str)
	}

	// 中间结果（图）构造，也就是用于回溯操作的对象
	for i := 0; i < len(table); i++ {
		obj = append(obj, table[i][0])
	}
	res = append(res, strings.Join(obj, ""))
	for i := pop(); i >= 0; i = pop() {
		l := len(table[i])
		if l > 1 {
			idx := calc[i] % l
			if idx == 0 && calc[i] > 0 {
				continue
			}
			rest := len(table) - len(obj)
			for j := 0; j < rest; j++ {
				lvl := len(obj)
				push(table[lvl][calc[lvl]%len(table[lvl])])
			}
			res = append(res, strings.Join(obj, ""))
		}
	}

	return res
}

func BracketExpansion2(s string) []string {

	var (
		res   []string
		table [][]string
	)

	for i := 0; i < len(s); i++ {
		if s[i] != '{' && s[i] != '}' && s[i] != ',' {
			table = append(table, []string{string(s[i])})
		} else if s[i] == '{' {
			collector := make([]string, 0, 4)
			for i++; s[i] != '}'; i++ {
				if s[i] == ',' {
					continue
				}
				collector = append(collector, string(s[i]))
			}
			isort.Sort(isort.StringSlice(collector))
			table = append(table, collector)
		}
	}

	calc := make([]int, len(table))
	var dfs func(lvl int, str string)
	// 使用深度优先遍历，完成对子树信息的获取，来替代使用栈的回溯实现
	dfs = func(lvl int, str string) {
		// 这里是多叉树的遍历，每一层的所有结点都要逐个访问
		for i := 0; lvl < len(table) && i < len(table[lvl]); i++ {
			s1 := table[lvl][calc[lvl]%len(table[lvl])]
			calc[lvl]++
			dfs(lvl+1, str+s1)
			if lvl == 0 {
				// 首层结点已经没有数据需要再访问，就停止此次遍历过程
				if calc[lvl]%len(table[lvl]) == 0 && calc[lvl] >= len(table[lvl]) {
					return
				}
			} else if lvl == len(table)-1 {
				// 最后一层结点直接拼接结果
				res = append(res, str+table[lvl][i])
			}
		}
	}
	dfs(0, "")
	return res
}
