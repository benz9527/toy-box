package hash

// References:
// https://colobu.com/2022/12/21/use-the-builtin-map-hasher/

import (
	"math"
)

func HashInt64(i int64) uint64 {
	return uint64(i)
}

func HashFloat64(f float64) uint64 {
	return math.Float64bits(f)
}

// HashString Using FNV-1a hash algorithm
func HashString(str string) uint64 {
	var hash uint64 = 14695981039346656037 // offset
	for i := 0; i < len(str); i++ {
		hash ^= uint64(str[i])
		hash *= 1099511628211
	}
	return hash
}
