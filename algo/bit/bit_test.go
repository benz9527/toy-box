package bit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoundupPowOf2(t *testing.T) {
	n := RoundupPowOf2(7)
	assert.Equal(t, RoundupPowOf2ByCeil(7), n)
	assert.Equal(t, RoundupPowOf2ByLoop(7), n)

	n = RoundupPowOf2(10)
	assert.Equal(t, RoundupPowOf2ByCeil(10), n)
	assert.Equal(t, RoundupPowOf2ByLoop(10), n)

	n = RoundupPowOf2(17)
	assert.Equal(t, RoundupPowOf2ByCeil(17), n)
	assert.Equal(t, RoundupPowOf2ByLoop(17), n)

	n = RoundupPowOf2(127)
	assert.Equal(t, RoundupPowOf2ByCeil(127), n)
	assert.Equal(t, RoundupPowOf2ByLoop(127), n)
}

func TestCeilPowOf2(t *testing.T) {
	n := CeilPowOf2(7)
	assert.Equal(t, uint8(3), n)

	n = CeilPowOf2(10)
	assert.Equal(t, uint8(4), n)

	n = CeilPowOf2(17)
	assert.Equal(t, uint8(5), n)
}

func TestOneBitsConvert(t *testing.T) {
	n := int8(-1)
	assert.Equal(t, uint64(255), convert[int8](n))

	n2 := int16(-1)
	assert.Equal(t, uint64(65535), convert[int16](n2))

}

func TestHammingWeight(t *testing.T) {
	n := 7
	assert.Equal(t, uint8(3), HammingWeightBySWAR[int](n))
	assert.Equal(t, uint8(3), HammingWeightBySWAR2[int](n))
	assert.Equal(t, uint8(3), HammingWeightBySWAR3[int](n))
	assert.Equal(t, uint8(3), HammingWeightByGroupCount[int](n))

	n2 := int64(0)
	assert.Equal(t, uint8(0), HammingWeightBySWAR[int64](n2))
	assert.Equal(t, uint8(0), HammingWeightBySWAR2[int64](n2))
	assert.Equal(t, uint8(0), HammingWeightBySWAR3[int64](n2))
	assert.Equal(t, uint8(0), HammingWeightByGroupCount[int64](n2))

	n3 := int8(-1)
	assert.Equal(t, uint8(8), HammingWeightBySWAR[int8](n3))
	assert.Equal(t, uint8(8), HammingWeightBySWAR2[int8](n3))
	assert.Equal(t, uint8(8), HammingWeightBySWAR3[int8](n3))
	assert.Equal(t, uint8(8), HammingWeightByGroupCount[int8](n3))

	n4 := int16(-1)
	assert.Equal(t, uint8(16), HammingWeightBySWAR[int16](n4))
	assert.Equal(t, uint8(16), HammingWeightBySWAR2[int16](n4))
	assert.Equal(t, uint8(16), HammingWeightBySWAR3[int16](n4))
	assert.Equal(t, uint8(16), HammingWeightByGroupCount[int16](n4))

	n5 := int32(-1)
	assert.Equal(t, uint8(32), HammingWeightBySWAR[int32](n5))
	assert.Equal(t, uint8(32), HammingWeightBySWAR2[int32](n5))
	assert.Equal(t, uint8(32), HammingWeightBySWAR3[int32](n5))
	assert.Equal(t, uint8(32), HammingWeightByGroupCount[int32](n5))

	n6 := int64(-1)
	assert.Equal(t, uint8(64), HammingWeightBySWAR[int64](n6))
	assert.Equal(t, uint8(64), HammingWeightBySWAR2[int64](n6))
	assert.Equal(t, uint8(64), HammingWeightBySWAR3[int64](n6))
	assert.Equal(t, uint8(64), HammingWeightByGroupCount[int64](n6))

}
