package bit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCeilPowerOfTwo(t *testing.T) {
	n := GetCeilPowerOfTwo(7)
	assert.Equal(t, uint64(8), n)

	n = GetCeilPowerOfTwo(10)
	assert.Equal(t, uint64(16), n)

	n = GetCeilPowerOfTwo(17)
	assert.Equal(t, uint64(32), n)
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
