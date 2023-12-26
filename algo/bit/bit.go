package bit

import (
	"unsafe"
)

func RoundupPowOf2(target uint64) uint64 {
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

// RoundupPowOf2ByLoop rounds up the target to the power of 2.
// Plain thinking.
func RoundupPowOf2ByLoop(target uint64) uint64 {
	var result uint64 = 1
	for result < target {
		result <<= 1
	}
	return result
}

// RoundupPowOf2ByCeil rounds up the target to the power of 2.
// Copy from linux kernel kfifo.
func RoundupPowOf2ByCeil(target uint64) uint64 {
	return 1 << CeilPowOf2(target)
}

// CeilPowOf2 get the ceil power of 2 of the target.
// Copy from linux kernel kfifo.
func CeilPowOf2(target uint64) uint8 {
	target--
	if target == 0 {
		return 0
	}
	var pos uint8 = 64
	if target&0xffffffff00000000 == 0 {
		target = target << 32
		pos -= 32
	}
	if target&0xffff000000000000 == 0 {
		target <<= 16
		pos -= 16
	}
	if target&0xff00000000000000 == 0 {
		target <<= 8
		pos -= 8
	}
	if target&0xf000000000000000 == 0 {
		target <<= 4
		pos -= 4
	}
	if target&0xc000000000000000 == 0 {
		target <<= 2
		pos -= 2
	}
	if target&0x8000000000000000 == 0 {
		pos -= 1
	}
	return pos
}

// IsPowOf2 checks if the target is power of 2.
// Copy from linux kernel kfifo.
func IsPowOf2(target uint64) bool {
	return target&(target-1) == 0
}

type Number interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// HammingWeight counts the number of 1 bit in a number.
type HammingWeight[T Number] func(n T) uint8
type BitCount[T Number] HammingWeight[T]

func convert[T Number](n T) uint64 {
	var target uint64
	switch unsafe.Sizeof(n) {
	case 1:
		target = uint64(*(*uint8)(unsafe.Pointer(&n)))
	case 2:
		target = uint64(*(*uint16)(unsafe.Pointer(&n)))
	case 4:
		target = uint64(*(*uint32)(unsafe.Pointer(&n)))
	case 8:
		target = *(*uint64)(unsafe.Pointer(&n))
	}
	return target
}

// variable-precision SWAR algorithm

// HammingWeightBySWAR counts the number of 1 bit by group statistics.
// Calculate the number of 1 bit by tree-like structure.
// Example:
// 0x55555555 = 01010101010101010101010101010101
// It will keep the odd bits of the original number 1 and keep the even bits of the original number 1.
// Every 2 binary bits represent the number of 1 in the corresponding binary bit of the original number.
// 0x33333333 = 00110011001100110011001100110011
// It will keep the right two bits of the sum of the previous step and keep the left two bits of the sum of the previous step.
// Every 4 binary bits represent the number of 1 in the corresponding binary bit of the original number.
// 0x0f0f0f0f = 00001111000011110000111100001111
// It will keep the right four bits of the sum of the previous step and keep the left four bits of the sum of the previous step.
// Every 8 binary bits represent the number of 1 in the corresponding binary bit of the original number.
// 0x00ff00ff = 00000000111111110000000011111111
// It will keep the right eight bits of the sum of the previous step and keep the left eight bits of the sum of the previous step.
// Every 16 binary bits represent the number of 1 in the corresponding binary bit of the original number.
// 0x0000ffff = 00000000000000001111111111111111
// It will keep the right sixteen bits of the sum of the previous step and keep the left sixteen bits of the sum of the previous step.
// Every 32 binary bits represent the number of 1 in the corresponding binary bit of the original number.
func HammingWeightBySWAR[T Number](n T) uint8 {
	_n := convert[T](n)
	_n = (_n & 0x5555555555555555) + ((_n >> 1) & 0x5555555555555555)
	_n = (_n & 0x3333333333333333) + ((_n >> 2) & 0x3333333333333333)
	_n = (_n & 0x0f0f0f0f0f0f0f0f) + ((_n >> 4) & 0x0f0f0f0f0f0f0f0f)
	_n = (_n & 0x00ff00ff00ff00ff) + ((_n >> 8) & 0x00ff00ff00ff00ff)
	_n = (_n & 0x0000ffff0000ffff) + ((_n >> 16) & 0x0000ffff0000ffff)
	_n = (_n & 0x00000000ffffffff) + ((_n >> 32) & 0x00000000ffffffff)
	return uint8(_n)
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
func HammingWeightBySWAR2[T Number](n T) uint8 {
	_n := convert[T](n)
	_n = (_n & 0x5555555555555555) + ((_n >> 1) & 0x5555555555555555)
	_n = (_n & 0x3333333333333333) + ((_n >> 2) & 0x3333333333333333)
	_n = (_n & 0x0f0f0f0f0f0f0f0f) + ((_n >> 4) & 0x0f0f0f0f0f0f0f0f)
	// 8 bits quick multiply
	// 0x01010101 = 00000001 00000001 00000001 00000001
	//            = 1 << 24 | 1 << 16 | 1 << 8 | 1 << 0
	// i * 0x01010101 = i << 24 + i << 16 + i << 8 + i << 0
	// Merge
	// (i * 0x01010101)>>24 = (i<<24)>>24 + (i<<16)>>24 + (i<<8)>>24 + (i<<0)>>24
	// Hamming Weight
	_n = ((_n * 0x0101010101010101) & ((1 << 64) - 1)) >> 56
	return uint8(_n)
}

// HammingWeightBySWAR3 counts the number of 1 bit by group statistics.
func HammingWeightBySWAR3[T Number](n T) uint8 {
	_n := convert[T](n)
	bits := func(num uint8) uint8 {
		remainder := num&0x5 + (num>>1)&0x5
		return remainder&0x3 + (remainder>>2)&0x3
	}
	res := uint8(0)
	for i := 0; i < 16; i++ {
		res += bits(uint8((_n >> (i * 4)) & 0xf))
	}
	return res
}

var (
	bitCount = [16]uint8{
		0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
	}
)

func HammingWeightByGroupCount[T Number](n T) uint8 {
	_n := convert[T](n)
	res := uint8(0)
	for i := 0; i < 16; i++ {
		res += bitCount[uint8((_n>>(i*4))&0xf)]
	}
	return res
}
