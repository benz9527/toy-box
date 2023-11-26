package assembly

import (
	"unsafe"

	"github.com/klauspost/cpuid/v2"
	"golang.org/x/exp/constraints"
)

// Experimental, high risk and unsafe.
// Not recommended to use in production.

var (
	// In debugging mode, set it as false
	hasERMS  = cpuid.CPU.Has(cpuid.ERMS)
	isX86_64 = cpuid.CPU.X64Level() != 0
)

type CanMove interface {
	constraints.Integer |
		constraints.Float |
		constraints.Complex |
		~string |
		uintptr
}

//go:linkname memmove runtime.memmove
func memmove(to, from unsafe.Pointer, n uintptr)

// Assembly definitions mapping function body automatically.

func can_movs(to, from unsafe.Pointer, n uintptr) (ok bool)
func copy_movsb(to, from unsafe.Pointer, n uintptr)
func copy_movsq(to, from unsafe.Pointer, n uintptr) (left, copied uintptr)

func FastCopy[T CanMove](dst, src []T) (n int) {
	n = min(len(src), len(dst))
	if n <= 0 {
		return
	}
	pDst := unsafe.Pointer(&dst[0])
	pSrc := unsafe.Pointer(&src[0])
	if pDst == pSrc {
		n = 0
		return
	}
	pN := uintptr(n)
	elementSize := unsafe.Sizeof(src[0])
	size := pN * elementSize
	if size > 15500 && (hasERMS || isX86_64) && can_movs(pDst, pSrc, pN) {
		if hasERMS {
			copy_movsb(pDst, pSrc, size)
		} else {
			left, copied := copy_movsq(pDst, pSrc, size)
			if left > 0 {
				memmove(unsafe.Pointer(&dst[copied]), unsafe.Pointer(&src[copied]), left*elementSize)
			}
		}
	} else {
		memmove(pDst, pSrc, size)
	}
	return
}

func FastCopyByMOVSB[T CanMove](dst, src []T) (n int) {
	n = min(len(src), len(dst))
	if n <= 0 {
		return
	}
	pDst := unsafe.Pointer(&dst[0])
	pSrc := unsafe.Pointer(&src[0])
	if pDst == pSrc {
		n = 0
		return
	}
	pN := uintptr(n)
	elementSize := unsafe.Sizeof(src[0])
	size := pN * elementSize
	if can_movs(pDst, pSrc, pN) {
		copy_movsb(pDst, pSrc, size)
	} else {
		memmove(pDst, pSrc, size)
	}
	return
}

func FastCopyByMOVSQ[T CanMove](dst, src []T) (n int) {
	n = min(len(src), len(dst))
	if n <= 0 {
		return
	}
	pDst := unsafe.Pointer(&dst[0])
	pSrc := unsafe.Pointer(&src[0])
	if pDst == pSrc {
		n = 0
		return
	}
	pN := uintptr(n)
	elementSize := unsafe.Sizeof(src[0])
	size := pN * elementSize
	if can_movs(pDst, pSrc, pN) {
		left, copied := copy_movsq(pDst, pSrc, size)
		if left > 0 {
			memmove(unsafe.Pointer(&dst[copied]), unsafe.Pointer(&src[copied]), left*elementSize)
		}
	} else {
		memmove(pDst, pSrc, size)
	}
	return
}
