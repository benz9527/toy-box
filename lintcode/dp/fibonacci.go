package dp

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	// 1. dp data form：
	var _dp = make([]int, n+1) // 第 i 个位置为当前的斐波那契数
	// 2. state transition equation:
	// _dp[i] = _dp[i-1] + _dp[i-2] // 这个就是状态转移方程，或者说递归方程
	// 3. initialization:
	_dp[0], _dp[1] = 0, 1 // 初始化
	// 4. iteration order:
	for i := 2; i <= n; i++ { // 迭代的方向
		_dp[i] = _dp[i-1] + _dp[i-2]
	}
	return _dp[n]
}

func FibonacciOptimize(n int) int {
	if n <= 1 {
		return n
	}
	// 1. dp data form：
	var _dp []int // 第 i 个位置为当前的斐波那契数
	// 2. state transition equation:
	// _dp[i] = _dp[i-1] + _dp[i-2]
	// 3. initialization:
	_dp = []int{0, 1}
	// 4. iteration order:
	for i := 2; i <= n; i++ {
		sum := _dp[1] + _dp[0]
		// 减少空间使用量
		_dp[0], _dp[1] = _dp[1], sum
	}
	return _dp[1]
}
