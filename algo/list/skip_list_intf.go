package list

import (
	"math"
	"math/bits"
)

// Ref
// paper:
// https://www.cl.cam.ac.uk/teaching/2005/Algorithms/skiplists.pdf
// github:
// classic: https://github.com/antirez/disque/blob/master/src/skiplist.h
// classic: https://github.com/antirez/disque/blob/master/src/skiplist.c
// https://github.com/liyue201/gostl
// https://github.com/chen3feng/stl4go
// test:
// https://github.com/chen3feng/skiplist-survey

const (
	ClassicSkipListMaxLevel    = 32   // 2^32 - 1 elements
	ClassicSkipListProbability = 0.25 // P = 1/4, a skip list node element has 1/4 probability to have a level
)

type SkipListNodeElement[E comparable] interface {
	GetObject() E
	GetVerticalBackward() SkipListNodeElement[E]
	SetVerticalBackward(backward SkipListNodeElement[E])
	GetLevels() []SkipListLevel[E]
	Free()
}

type SkipListLevel[E comparable] interface {
	GetSpan() int64
	SetSpan(span int64)
	GetHorizontalForward() SkipListNodeElement[E]
	SetHorizontalForward(forward SkipListNodeElement[E])
}

type SkipList[E comparable] interface {
	GetLevel() int
	Len() int64
	Insert(v E) SkipListNodeElement[E]
	Remove(v E) SkipListNodeElement[E]
	Find(v E) SkipListNodeElement[E]
	PopHead() E
	PopTail() E
	Free()
	ForEach(fn func(idx int64, v E))
}

// LessThan is the compare function.
type compareTo[E comparable] func(a, b E) int

// SkipListRandomLevel is the skip list level element.
// Dynamic level calculation.
func SkipListRandomLevel(random func() uint64, maxLevel int) int {
	// goland math random (math.Float64()) contains global mutex lock
	// Ref
	// https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/math/rand/rand.go
	// https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/math/bits/bits.go
	// 1. Avoid to use global mutex lock
	// 2. Avoid to generate random number each time
	total := uint64(1)<<maxLevel - 1 // maxLevel => n, 2^n -1, there will be 2^n-1 elements in the skip list
	rest := random() % total
	// Bits right shift equals to manipulate a high level bit
	// Calculate the minimum bits of the random number
	level := maxLevel - bits.Len64(rest) + 1
	return level
}

func MaxLevels(totalElements int64, P float64) int {
	// Ref https://www.cl.cam.ac.uk/teaching/2005/Algorithms/skiplists.pdf
	// MaxLevels = log(1/P) * log(totalElements)
	// P = 1/4, totalElements = 2^32 - 1
	return int(math.Ceil(math.Log(1/P) * math.Log(float64(totalElements))))
}
