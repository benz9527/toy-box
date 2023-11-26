package assembly

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test_can_movs(t *testing.T) {
	ori := make([]byte, 32+16)
	from := ori[1 : 32+1]
	to := ori[16 : 32+16]
	t.Log(can_movs(unsafe.Pointer(&from[0]), unsafe.Pointer(&to[0]), 32))
}

func TestFastCopy_byte(t *testing.T) {
	size := 32768
	src := make([]byte, size)
	dst := make([]byte, size)
	for i := 0; i < size; i++ {
		src[i] = byte(i)
	}
	t.Log(hasERMS, isX86_64, FastCopy[byte](dst, src))
	for i := 0; i < len(dst); i++ {
		if dst[i] != byte(i) {
			t.Fatal("dst is not equal to src, fast move failed")
		}
	}
}

func TestFastCopy_float32(t *testing.T) { // 4 bytes
	size := 32768
	src := make([]float32, size)
	dst := make([]float32, size)
	for i := 0; i < size; i++ {
		src[i] = float32(i)
	}
	t.Log(hasERMS, isX86_64, FastCopy[float32](dst, src))
	for i := 0; i < size; i++ {
		if dst[i] != float32(i) {
			t.Fatal("dst is not equal to src, fast move failed")
		}
	}
}

func TestFastCopy_int(t *testing.T) { // 8 bytes
	size := 32768
	src := make([]int, size)
	dst := make([]int, size)
	for i := 0; i < size; i++ {
		src[i] = i
	}
	t.Log(hasERMS, isX86_64, FastCopy[int](dst, src))
	for i := 0; i < size; i++ {
		if dst[i] != i {
			t.Fatal("dst is not equal to src, fast move failed")
		}
	}
}

func TestFastCopy_slice(t *testing.T) {
	src := make([]int, 100)
	dst := make([]int, 1000)
	for i := 0; i < len(src); i++ {
		src[i] = i
	}
	var n int
	for i := 0; i < 10; i++ {
		n += FastCopy[int](dst[n:], src)
	}
	t.Log(hasERMS, isX86_64, dst)
}

func BenchmarkFastCopy_int(b *testing.B) {
	size := 32768
	src := make([]int, size)
	dst := make([]int, size)
	for i := 0; i < size; i++ {
		src[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FastCopy[int](dst, src)
	}
}

func BenchmarkFastCopy_int_ByGoCopy(b *testing.B) {
	size := 32768
	src := make([]int, size)
	dst := make([]int, size)
	for i := 0; i < size; i++ {
		src[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(dst, src) // Slower than FastCopy
	}
}

func benchmarkSizes(b *testing.B, sizes []int, fn func(b *testing.B, size int)) {
	for _, n := range sizes {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			b.SetBytes(int64(n))
			fn(b, n)
		})
	}
}

func BenchmarkFastCopyByMOVSBOutput(b *testing.B) {
	largeBuffers := []int{
		1024, 2048, 4096, 8192, 15500, 16000, 16384, 25000, 32768, 65536,
	}
	benchmarkSizes(b, largeBuffers, func(b *testing.B, size int) {
		src := make([]byte, size)
		dst := make([]byte, size)
		for i := 0; i < b.N; i++ {
			FastCopyByMOVSB[byte](dst, src)
		}
	})
}

func BenchmarkFastCopyByMOVSQOutput(b *testing.B) {
	largeBuffers := []int{
		1024, 2048, 4096, 8192, 15500, 16000, 16384, 25000, 32768, 65536,
	}
	benchmarkSizes(b, largeBuffers, func(b *testing.B, size int) {
		src := make([]byte, size)
		dst := make([]byte, size)
		for i := 0; i < b.N; i++ {
			FastCopyByMOVSQ[byte](dst, src)
		}
	})
}
