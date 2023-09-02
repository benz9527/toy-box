package b

import (
	"strconv"
	"strings"
)

// 在一个博客网站上，每篇博客都有评论。
// 每一条评论都是一个非空英文字母字符串。
// 评论具有树状结构，除了根评论外，每个评论都有一个父评论。
// 当评论保存时，使用以下格式：
// 首先是评论的内容；
// 然后是回复当前评论的数量。
// 最后是当前评论的所有了评论。(子评论使用相同的格式嵌套存储)
// 所有元素之间都用单个逗号分隔。
// 第一条评论是"helo,2,ok,0,bye,0"，第二条评论是"test,0"，第三条评论是"one,1,two,1,a,0"。
// 所有评论被保存成"hello,2,ok,0.bye,0,test,0,one,1,two,1,a,0"。
// 对于上述格式的评论，请以另外一种格式打印：
// 首先打印评论嵌套的最大深度。
// 然后是打印n行，第 i (1 ≤ i ≤ n) 行对应于嵌套级别为 i 的评论 (根评论的嵌套级别为1)。
// 对于第 i 行，嵌套级别为的评论按照它们出现的顺序打印，用空格分隔开。

type myComment struct {
	comment string
	num     int
}

func CommentsTransfer(comments string) (int, []string) {
	tree := make(map[int]string, 16)
	input := strings.Split(comments, ",")
	newComments := make([]myComment, 0, len(input)/2)
	for i := 0; i < len(input); i += 2 {
		num, _ := strconv.Atoi(input[i+1])
		newComments = append(newComments, myComment{comment: input[i], num: num})
	}

	stack := make([]*myComment, 0, 16)
	pop := func() *myComment {
		res := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return res
	}
	push := func(c *myComment) {
		stack = append(stack, c)
	}
	peek := func() *myComment {
		return stack[len(stack)-1]
	}
	isEmpty := func() bool {
		return len(stack) == 0
	}
	for len(newComments) > 0 {
		push(&newComments[0]) // root
		newComments = newComments[1:]
		levels := 1
		for !isEmpty() {
			if peek().num == 0 {
				c := pop()
				if !isEmpty() {
					peek().num--
				}
				if _, ok := tree[levels]; !ok {
					tree[levels] = c.comment
				} else {
					tree[levels] += " " + c.comment
				}
				levels = len(stack)
				continue
			}
			push(&newComments[0])
			newComments = newComments[1:]
			levels = len(stack)
		}
	}

	levels := len(tree)
	ans := make([]string, 0, levels)
	for i := 1; i <= levels; i++ {
		ans = append(ans, tree[i])
	}

	return levels, ans
}
