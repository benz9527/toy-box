package arr

// https://leetcode.cn/problems/maximum-xor-of-two-numbers-in-an-array/description/

func FindMaximumXOR(nums []int) int {
	ans := 0
	root := &trie{}
	maximum := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	for i := 1; i < len(nums); i++ {
		root.add(nums[i-1]) // 每次遍历存入上一次的数据，最后 0,(i-1)的数值都会在字典树里面
		ans = maximum(ans, root.getXOR(nums[i]))
	}
	return ans
}

const highBits = 30

type trie struct {
	left, right *trie // 0, 1 node
}

func (t *trie) add(num int) {
	cur := t
	for i := highBits; i >= 0; i-- {
		ch := num >> i & 0x1
		if ch == 0x1 {
			if cur.right == nil {
				cur.right = &trie{}
			}
			cur = cur.right
		} else {
			if cur.left == nil {
				cur.left = &trie{}
			}
			cur = cur.left
		}
	}
}

func (t *trie) getXOR(num int) int {
	ans := 0
	cur := t
	for i := highBits; i >= 0; i-- {
		ch := num >> i & 0x1
		if ch == 0x1 {
			// 1 ^ 0 = 1，之后需要通过位移运算升位
			if cur.left != nil {
				ans = ans<<1 + 1
				cur = cur.left
			} else {
				ans = ans << 1
				cur = cur.right
			}
		} else {
			// 0 ^ 1 = 1
			if cur.right != nil {
				ans = ans<<1 + 1
				cur = cur.right
			} else {
				ans = ans << 1
				cur = cur.left
			}
		}
	}
	return ans
}
