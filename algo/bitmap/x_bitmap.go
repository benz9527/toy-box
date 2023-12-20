package bitmap

import "bytes"

const (
	maxBitMapSize = 1 << 32
)

type x32Bitmap struct {
	bits []byte
	size uint64
}

func NewX32Bitmap(size uint64) Bitmap {
	if size <= 0 || size > maxBitMapSize {
		size = maxBitMapSize
	}
	if remainder := size & 0x07; remainder != 0 {
		size = size + (8 - remainder)
	}
	return &x32Bitmap{
		bits: make([]byte, size>>3),
		size: size - 1,
	}
}

func (bm *x32Bitmap) SetBit(offset uint64, one bool) bool {
	idx, pos := offset>>3, offset&0x07
	if bm.size < offset {
		return false
	}
	if !one {
		bm.bits[idx] &= ^(1 << pos) // &^=
	} else {
		bm.bits[idx] |= 1 << pos
	}
	return true
}

func (bm *x32Bitmap) GetBit(offset uint64) bool {
	idx, pos := offset>>3, offset&0x07
	if bm.size < offset {
		return false
	}
	return bm.bits[idx]>>pos != 0
}

func (bm *x32Bitmap) GetBits() []byte {
	return bm.bits
}

func (bm *x32Bitmap) EqualTo(that Bitmap) bool {
	return bytes.Compare(bm.GetBits(), that.GetBits()) == 0
}
