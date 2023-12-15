package api

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"unsafe"
)

type myInterface interface {
	ToString() string
}

type myStruct struct {
	field1 int
}

func (s *myStruct) ToString() string {
	builder := &strings.Builder{}
	builder.WriteString("myStruct{")
	builder.WriteString("field1=")
	builder.Write([]byte{byte(s.field1)})
	builder.WriteString("}")
	return builder.String()
}

type myWrapper struct {
	ms unsafe.Pointer
}

type myInterface2[E comparable] interface {
	GetE() E
}

type myStruct2[E comparable] struct {
	field1 E
}

func (s *myStruct2[E]) GetE() E {
	return s.field1
}

type myWrapper2[E comparable] struct {
	mi2 unsafe.Pointer // E
}

func TestUnsafePointerConvert(t *testing.T) {
	mw := myWrapper{}
	s := myStruct{field1: 1}
	atomic.StorePointer(&mw.ms, unsafe.Pointer(&s))
	ms := (*myStruct)(atomic.LoadPointer(&mw.ms))
	t.Log(ms.ToString())

	var s1 myInterface = &myStruct{field1: 1}
	atomic.StorePointer(&mw.ms, unsafe.Pointer(&s1))
	ms1 := (*myInterface)(atomic.LoadPointer(&mw.ms))
	t.Log((*ms1).ToString())

	mw2 := myWrapper2[int]{}
	i := 1
	atomic.StorePointer(&mw2.mi2, unsafe.Pointer(&i))
	ms2Field := *(*int)(atomic.LoadPointer(&mw2.mi2))
	t.Log(ms2Field)
}

func TestSliceHeaderAddrManipulate(t *testing.T) {
	s := make([]int, 10)

	_1stPtrOfS := unsafe.SliceData(s)
	*_1stPtrOfS = 1
	assert.Equal(t, 1, s[0])

	ps := unsafe.Slice(&s[0], 10)
	for i := 0; i < 10; i++ {
		ps[i] = i + 5
	}
	expected := []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	assert.Equal(t, expected, s)

	// 1.21 deprecated
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	_1stEptr := unsafe.Pointer(sh.Data)
	for i := 0; i < 10; i++ {
		*(*int)(unsafe.Pointer(uintptr(_1stEptr) + uintptr(i)*unsafe.Sizeof(int(0)))) = i + 10
	}
	expected = []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	assert.Equal(t, expected, s)
}

func TestFastSliceConvert(t *testing.T) {
	string2Byets := func(str string) []byte {
		if str == "" {
			return nil
		}
		// cap is equal to len of str
		return unsafe.Slice(unsafe.StringData(str), len(str))
	}
	string2BytesOld := func(str string) []byte {
		if str == "" {
			return nil
		}
		// cap is too large, more than len of str
		return *(*[]byte)(unsafe.Pointer(&str))
	}
	bytes2String := func(bytes []byte) string {
		if bytes == nil {
			return ""
		}
		return unsafe.String(unsafe.SliceData(bytes), len(bytes))
	}
	bytes2StringOld := func(bytes []byte) string {
		if bytes == nil {
			return ""
		}
		return *(*string)(unsafe.Pointer(&bytes))
	}

	str := "hello world"
	bytes1 := string2Byets(str)
	assert.Equal(t, str, bytes2String(bytes1))
	bytes2 := string2BytesOld(str)
	assert.Equal(t, str, bytes2StringOld(bytes2))

	bytes3 := []byte(str)
	assert.Equal(t, str, bytes2String(bytes3))
	assert.Equal(t, str, bytes2StringOld(bytes3))
}
