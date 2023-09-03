package b

import isort "sort"

// 二分法

// 一个设备由N种类型元器件组成(每种类型元器件只需要一个，类型type编号从0~N-1)，
// 每个元器件均有可靠性属性reliability，可靠性越高的器件其价格price越贵。
// 而设备的可靠性由组成设备的所有器件中可靠性最低的器件决定。
// 给定预算S，购买N种元器件( 每种类型元器件都需要购买一个)，在不超过预算的情况下，请给出能够组成的设备的最大可靠性。
// S 总预算 N 元器件类型
// total 元器件总数
// type 元器件类型编号 reliability 元器件的可靠性 price 元器件的价格

type component struct {
	kind, price, reliability int
}
type componentSlice []component

func (s componentSlice) Less(i, j int) bool {
	return s[i].reliability < s[j].reliability
}
func (s componentSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s componentSlice) Len() int {
	return len(s)
}

func BuyMachines(totalSum, totalComponents int, components []component) int {
	reliabilities := make([]int, 0, len(components))
	groupedComponents := map[int][]component{}
	for i := 0; i < len(components); i++ {
		reliabilities = append(reliabilities, components[i].reliability)
		_, ok := groupedComponents[components[i].kind]
		if !ok {
			groupedComponents[components[i].kind] = []component{components[i]}
		} else {
			groupedComponents[components[i].kind] = append(groupedComponents[components[i].kind], components[i])
		}
	}
	isort.Ints(reliabilities)
	for _, v := range groupedComponents {
		isort.Sort(componentSlice(v))
	}

	// 根据 reliability 约高，price 也越高的同步关系，可以用二分查找
	binarySearch := func(components []component, maxReliability int) int {
		// 返回索引
		lo, hi := 0, len(components)-1
		for lo <= hi {
			midd := (lo + hi) >> 1
			if components[midd].reliability > maxReliability {
				hi = midd - 1
			} else if components[midd].reliability < maxReliability {
				lo = midd + 1
			} else {
				return midd
			}
		}
		// 插空位置
		return -lo - 1
	}
	isSatisfy := func(gComponents map[int][]component, maxReliability, totalComponents int) bool {
		sum := 0
		for _, v := range gComponents {
			idx := binarySearch(v, maxReliability)
			if idx < 0 {
				idx = -idx - 1
			}
			if idx == len(v) {
				// 所有都小
				return false
			}
			sum += v[idx].price
		}
		return sum <= totalSum
	}

	l, h, ans := 0, len(components)-1, -1
	for l <= h {
		mid := (l + h) >> 1
		if isSatisfy(groupedComponents, reliabilities[mid], totalComponents) {
			ans = reliabilities[mid]
			l = mid + 1
		} else {
			h = mid - 1
		}
	}
	return ans
}

// 新来的老师给班里的同学排一个队。
// 每个学生有一个影力值。
// 一些学生是刺头，不会听老师的话，自己选位置，非刺头同学在剩下的位置按照能力值从小到大排。
// 对于非刺头同学，如果发现他前面有能力值比自己高的同学，他不满程度就增加，增加的数量等于前面能力值比他大的同学的个数。
// 刺头不会产生不满。
// 如果整个班级累计的不满程度超过k，那么老师就没有办法教这个班级了。
const (
	teachable    = 0
	disteachable = 1
)

func AngryStudentsAreTeachable(allStudents, badStudentIdxs []int, tolerance int) int {
	ans := teachable
	badStuSet := map[int]struct{}{}
	for _, idx := range badStudentIdxs {
		badStuSet[idx] = struct{}{}
	}
	badStudentAbis := make([]int, 0, 8)
	angry := 0
	_getLastIdx := func(abi int) int {
		l, r := 0, len(badStudentAbis)-1
		for l <= r {
			mid := (l + r) >> 1
			badAbi := badStudentAbis[mid]
			if abi < badAbi {
				r = mid - 1
			} else if abi > badAbi {
				l = mid + 1
			} else {
				if mid == len(badStudentAbis)-1 || badAbi-badStudentAbis[mid+1] != 0 {
					return mid
				}
				l = mid + 1 // unique
			}
		}
		return -l - 1
	}
	for i := 0; i < len(allStudents); i++ {
		abi := allStudents[i]
		idx := _getLastIdx(abi)

		if idx < 0 {
			idx = -idx - 1
		} else {
			idx += 1
		}

		if _, ok := badStuSet[i]; ok {
			if len(badStudentAbis) == 0 || len(badStudentAbis) > 0 && len(badStudentIdxs) < idx {
				badStudentAbis = append(badStudentAbis, abi)
			} else if len(badStudentAbis) > 0 && len(badStudentIdxs) >= idx {
				badStudentAbis = append(badStudentAbis[:idx], append([]int{abi}, badStudentAbis[idx:]...)...)
			}
		} else {
			angry += len(badStudentAbis) - idx
		}
	}
	if angry > tolerance {
		ans = disteachable
	}
	return ans
}
