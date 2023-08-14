package feature

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Go1.21.0
// func max[T cmp.Ordered](x T, y ...T) T
// func min[T cmp.Ordered](x T, y ...T) T

func TestGo1210DefaultMaxFunc(t *testing.T) {
	var m1, m2, m3 int64 = 1, 2, 3
	res := max(m1, m2, m3)
	assert.Equal(t, m3, res)
	typ := fmt.Sprintf("%T", res)
	assert.Equal(t, "int64", typ)

	var f1, f2, f3 float32 = 1.0, 2.0, 3.01
	fres := max(f1, f2, f3)
	assert.Equal(t, f3, fres)
	typ = fmt.Sprintf("%T", fres)
	assert.Equal(t, "float32", typ)

	// 字典序比较
	var s1, s2, s3 = "1a", "1b", "1c"
	sres := max(s1, s2, s3)
	assert.Equal(t, s3, sres)
	typ = fmt.Sprintf("%T", sres)
	assert.Equal(t, "string", typ)
}
