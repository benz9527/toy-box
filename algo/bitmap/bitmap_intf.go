package bitmap

type Bitmap interface {
	SetBit(offset uint64, one bool) bool
	GetBit(offset uint64) bool
}
