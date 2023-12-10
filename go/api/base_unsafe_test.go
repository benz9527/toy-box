package api

import (
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
