package b

import isort "sort"

// Solo 和 koko 是两兄弟，妈妈给了他们一大堆积木，每块积木上都有自己的重量。
// 现在他们想要将这些积木分成两堆。哥哥 Solo 负责分配，弟弟 koko 要求两个人
// 获得的积木总重量 “相等”（根据 Koko 的逻辑），个数可以不同，不然就会哭，
// 但 koko 只会先将两个数转成二进制再进行加法，而且总会忘记进位（每个进位都忘记）。
// 如当 25（11101）加 11（01011）时，koko 得到的计算结果是 18（10010）
// Solo 想要尽可能使自己得到的积木总重量最大，且不让 koko 哭。

func BricksEqualDivision(bricks []int) (int, bool) {
	// 所谓忘记进位，其实就是按位异或
	// 异或中的平均分配就是参与运算的是两个相等的数，
	// 更进一步说就是所有的元素参与运算之后，结果为 0
	// 要使的这种“公平均分”其中一堆达到最大，只要所有
	// 元素异或结果为 0，那么常规加法求和中的总和减去
	// 最小的元素值即可
	sumAll := func(arr []int) int {
		total := 0
		for _, e := range arr {
			total += e
		}
		return total
	}
	xorAll := func(arr []int) bool {
		// 000 ^ 010 = 010
		res := 0
		for _, e := range arr {
			res ^= e
		}
		return res == 0
	}
	isort.Ints(bricks)
	minBrick := bricks[0]
	sum := sumAll(bricks)
	if xor := xorAll(bricks); !xor {
		return -1, false
	}
	return sum - minBrick, true
}
