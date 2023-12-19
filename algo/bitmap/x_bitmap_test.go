package bitmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewX32Bitmap(t *testing.T) {
	bm := NewX32Bitmap(10)
	originalOffsets := []uint64{9, 5, 7, 3, 2, 8, 1}
	expectedOffsets := []uint64{1, 2, 3, 5, 7, 8, 9}
	for _, offset := range originalOffsets {
		bm.SetBit(offset, true)
	}
	for _, offset := range expectedOffsets {
		assert.True(t, bm.GetBit(offset))
	}
}
