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

func TestHammingWeight(t *testing.T) {
	n := 7
	assert.Equal(t, int64(3), HammingWeightBySWAR[int](n))
	assert.Equal(t, int64(3), HammingWeightBySWAR2[int](n))
	assert.Equal(t, int64(3), HammingWeightBySWAR3[int](n))
	assert.Equal(t, int64(3), HammingWeightByGroupCount[int](n))

	n2 := int64(-1)
	assert.Equal(t, int64(64), HammingWeightBySWAR[int64](n2))
	assert.Equal(t, int64(64), HammingWeightBySWAR2[int64](n2))
	assert.Equal(t, int64(64), HammingWeightBySWAR3[int64](n2))
	assert.Equal(t, int64(64), HammingWeightByGroupCount[int64](n2))
}
