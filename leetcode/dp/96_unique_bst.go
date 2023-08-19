package dp

// https://leetcode.cn/problems/unique-binary-search-trees/

func NumTrees(n int) int {
	if n == 0 {
		return 0
	}
	if n < 2 {
		return 1
	}
	_dp := make([]int, n+1)
	_dp[0] = 1 // 辅助判断，表示所在的子树没有结点变换，但是可取方案数系数为 1
	_dp[1] = 1 // 辅助判断，表示所在的子树没有结点变换，但是可取方案数系数为 1
	// _dp[2] 展开逻辑为 _dp[0]*_dp[1] + _dp[1]*dp[0]
	// 详细解释为
	// 1. 取 1 为根，只会有单结点右子树（bst 大于根的节点在右边，否则在左边）
	// 2. 取 2 为根，只会有单结点左子树
	_dp[2] = 2

	for i := 3; i <= n; i++ {
		for j := 0; j < i; j++ {
			// 不断乱换结点做根，左右子树的变换乘积之和即为 bst 的数量
			// 因为结点值是递增的，左右子树结点数量的变化在遍历中逐渐可知的
			// 比如一颗子树中只要是3个结点，它的 bst 子树数量就是 5
			_dp[i] += _dp[j] * _dp[i-1-j]
		}
	}
	return _dp[n]
}
