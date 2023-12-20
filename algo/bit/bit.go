package bit

func GetCeilPowerOfTwo(target uint64) uint64 {
	target--
	target |= target >> 1
	target |= target >> 2
	target |= target >> 4
	target |= target >> 8
	target |= target >> 16
	target |= target >> 32
	target++
	return target
}

type OneBits interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// HammingWeight counts the number of 1 bit in a number.
type HammingWeight[T OneBits] func(n T) int64

// HammingWeightBySWAR counts the number of 1 bit by group statistics.
func HammingWeightBySWAR[T OneBits](n T) int64 {
	_n := int64(n)
	// 使用树状的方式进行计算
	// 0x55555555 = 01010101010101010101010101010101
	// 保留原数字的奇数位的 1 + 保留原数字的偶数位的 1
	// 每两个二进制位表示原数字对应的二进制位的1的数量
	_n = (_n & 0x5555555555555555) + ((_n >> 1) & 0x5555555555555555)
	// 0x33333333 = 00110011001100110011001100110011
	// 保留上一步求出的和的右边两个位 + 保留上一步求出的和的左边两个位
	// 每四个二进制位表示原数字对应二进制位的 1 的数量
	_n = (_n & 0x3333333333333333) + ((_n >> 2) & 0x3333333333333333)
	// 0x0f0f0f0f = 00001111000011110000111100001111
	// 保留上一步求出的右边四个位 + 保留上一步求出的和的左边四个位
	// 每八个二进制位表示原数字对应二进制位的1的数量
	_n = (_n & 0x0f0f0f0f0f0f0f0f) + ((_n >> 4) & 0x0f0f0f0f0f0f0f0f)
	// 0x00ff00ff = 00000000111111110000000011111111
	// 保留上一步求出的右边八个位 + 保留上一步求出的和的左边八个位
	// 每十六个二进制位表示原数字对应二进制位的1的数量
	_n = (_n & 0x00ff00ff00ff00ff) + ((_n >> 8) & 0x00ff00ff00ff00ff)
	// 0x0000ffff = 00000000000000001111111111111111
	// 保留上一步求出的右边十六个位 + 保留上一步求出的和的左边十六个位
	// 每三十二个二进制位表示原数字对应二进制位的1的数量
	_n = (_n & 0x0000ffff0000ffff) + ((_n >> 16) & 0x0000ffff0000ffff)
	_n = (_n & 0x00000000ffffffff) + ((_n >> 32) & 0x00000000ffffffff)
	return _n
}

// HammingWeightBySWAR2 counts the number of 1 bit by group statistics.
// Example:
// 7 = (0111)2
// step 1:
// 0x7 & 0x55555555 = 0x5
// 0x7 >> 1 = 0x3, 0x3 & 0x55555555 = 0x1
// 0x5 + 0x1 = 0x6
// step 2:
// 0x6 & 0x33333333 = 0x2
// 0x6 >> 2 = 0x1, 0x1 & 0x33333333 = 0x1
// 0x2 + 0x1 = 0x3
// step 3:
// 0x3 & 0x0f0f0f0f = 0x3
// 0x3 >> 4 = 0x0, 0x0 & 0x0f0f0f0f = 0x0
// 0x3 + 0x0 = 0x3
// step 4:
// 0x3 * 0x01010101 = 0x03030303
// 0x03030303 & 0x3fffffff = 0x03030303
// 0x03030303 >> 24 = 0x3
func HammingWeightBySWAR2[T OneBits](n T) int64 {
	_n := int64(n)
	// 以 32 位为例
	// 0x55555555 = 01010101010101010101010101010101
	// 保留原数字的奇数位的 1 + 保留原数字的偶数位的 1
	// 每两个二进制位表示原数字对应的二进制位的1的数量
	_n = (_n & 0x5555555555555555) + ((_n >> 1) & 0x5555555555555555)
	// 0x33333333 = 00110011001100110011001100110011
	// 保留上一步求出的和的右边两个位 + 保留上一步求出的和的左边两个位
	// 每四个二进制位表示原数字对应二进制位的 1 的数量
	_n = (_n & 0x3333333333333333) + ((_n >> 2) & 0x3333333333333333)
	// 0x0f0f0f0f = 00001111000011110000111100001111
	// 保留上一步求出的右边四个位 + 保留上一步求出的和的左边四个位
	// 每八个二进制位表示原数字对应二进制位的1的数量
	_n = (_n & 0x0f0f0f0f0f0f0f0f) + ((_n >> 4) & 0x0f0f0f0f0f0f0f0f)
	// 32 bit reach here will be overed by ((_n * 0x01010101) & ((1 << 32) - 1)) >> 24
	_n = (_n & 0x00ff00ff00ff00ff) + ((_n >> 8) & 0x00ff00ff00ff00ff)
	_n = (_n & 0x0000ffff0000ffff) + ((_n >> 16) & 0x0000ffff0000ffff)
	_n = (_n & 0x00000000ffffffff) + ((_n >> 32) & 0x00000000ffffffff)
	// Hamming Weight
	_n = ((_n * 0x01010101001010101) & ((1 << 32) - 1)) >> 24
	return _n
}

// HammingWeightBySWAR3 counts the number of 1 bit by group statistics.
func HammingWeightBySWAR3[T OneBits](n T) int64 {
	_n := int64(n)
	bits := func(num int64) int64 {
		remainder := num&0x5 + (num>>1)&0x5
		return remainder&0x3 + (remainder>>2)&0x3
	}
	res := int64(0)
	for i := 0; i < 16; i++ {
		res += bits((_n >> (i * 4)) & 0xf)
	}
	return res
}

var (
	bitCount = [16]uint8{
		0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
	}
)

func HammingWeightByGroupCount[T OneBits](n T) int64 {
	res := int64(0)
	_n := int64(n)
	for i := 0; i < 16; i++ {
		res += int64(bitCount[(_n>>(i*4))&0xf])
	}
	return res
}
