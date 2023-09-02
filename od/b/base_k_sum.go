package b

import isort "sort"

type bigIntSlice []int64

func (s bigIntSlice) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s bigIntSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s bigIntSlice) Len() int {
	return len(s)
}

func KSum(nums []int64, k int, target int64) int {
	if k > len(nums) {
		return 0
	}
	isort.Sort(bigIntSlice(nums))
	return kSum(nums, k, 0, 0, target, 0)
}

func kSum(nums []int64, k, start, count int, target, sum int64) int {
	if k < 2 {
		return count
	}

	if k == 2 {
		return _2Sum(nums, start, count, target, sum)
	}

	for i := start; i < len(nums)-k; i++ {
		if nums[i] > 0 && sum+nums[i] > target {
			break
		}
		if i > start && nums[i] == nums[i-1] {
			continue
		}
		count = kSum(nums, k-1, i+1, count, target, sum+nums[i])
	}
	return count
}

func _2Sum(nums []int64, start, count int, target, sum int64) int {
	l, r := start, len(nums)-1
	for l < r {
		curSum := sum + nums[l] + nums[r]
		if target < curSum {
			r--
		} else if target > curSum {
			l++
		} else {
			count++
			for l+1 < r && nums[l] == nums[l+1] {
				l++
			}
			for r-1 > l && nums[r] == nums[r-1] {
				r--
			}
			l++
			r--
		}
	}
	return count
}

// 跳房子，也叫跳飞机，是一种世界性的儿童游戏。
// 游戏参与者需要分多个回合按顺序跳到第1格直到房子的最后一格，然后获得一次选房子的机会，直到所有房子被选完，房子最多的人获胜。
// 跳房子的过程中，如果有踩线等违规行为，会结束当前回合，甚至可能倒退几步。
// 假设房子的总格数是count，小红每回合可能连续跳的步数都放在数组steps中，请问数组中是否有一种步数的组合，
// 可以让小红三个回合跳到最后一格?
// 如果有，请输出索引和最小的步数组合（数据保证索引和最小的步数组合是唯一的）。
// 注意：数组中的步数可以重复，但数组中的元素不能重复使用。

type myStep struct {
	val, idx int
}
type mySteps []myStep

func (m mySteps) Len() int {
	return len(m)
}

func (m mySteps) Less(i, j int) bool {
	if m[i].val == m[j].val {
		return m[i].idx < m[j].idx
	}
	return m[i].val < m[j].val
}

func (m mySteps) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func JumpHousesII(steps []int, count int) []int {
	newSteps := make([]myStep, 0, len(steps))
	for i := 0; i < len(steps); i++ {
		newSteps = append(newSteps, myStep{idx: i, val: steps[i]})
	}
	ansSteps := make([]int, 0, 3)
	isort.Sort(mySteps(newSteps))
	minIdxSum := 10000 * 10000
	_2sum := func(start int, upperStep myStep) {
		l, r := start, len(newSteps)-1
		valSum := upperStep.val
		idxSum := upperStep.idx
		for l < r {
			// 这个也可以不需要，也是剪枝的操作
			// L,R 指针指向值的目标和为 (总目标值 - 前面固定元素对应的值之和)，而 L指针指向的值必然小于等于 R 指针指向的值，
			// 减少无用步骤
			threshold := (count - valSum) / 2
			if newSteps[l].val > threshold || newSteps[r].val < threshold {
				break
			}

			// 获取之前固定元素和当前活动元素 L R 之间的总和
			curSum := valSum + newSteps[l].val + newSteps[r].val
			if count > curSum {
				l++
			} else if count < curSum {
				r--
			} else { // 达到目标和
				// 值相等的情况下，需要剪枝
				for r-1 > l && newSteps[r].val == newSteps[r-1].val {
					r--
				}

				// 处理 index sum
				curIdxSum := idxSum + newSteps[l].idx + newSteps[r].idx
				if curIdxSum < minIdxSum {
					minIdxSum = curIdxSum
					tmp := []int{upperStep.idx, newSteps[l].idx, newSteps[r].idx}
					isort.Ints(tmp)
					ansSteps = []int{steps[tmp[0]], steps[tmp[1]], steps[tmp[2]]}
				}

				// 值相等的情况下，需要剪枝
				for l+1 < r && newSteps[l].val == newSteps[l+1].val {
					l++
				}
				l++
				r--
			}
		}
	}
	var _ksum func(k, start int, step ...myStep)
	_ksum = func(k, start int, step ...myStep) {
		if k < 2 {
			return
		}

		if k == 2 {
			_2sum(start, step[0])
			return
		}

		for i := start; i < len(steps)-k; i++ {
			if newSteps[i].val > 0 && newSteps[i].val > count {
				break
			}
			if i > start && newSteps[i].val == newSteps[i-1].val {
				continue
			}
			_ksum(k-1, i+1, newSteps[i])
		}
	}
	_ksum(3, 0)
	return ansSteps
}
