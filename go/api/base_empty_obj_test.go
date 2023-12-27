package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

type obj interface {
	abc()
}

type obj1 struct {
	s string
}

func (o *obj1) abc() {

}

type obj2 struct {
}

func (o *obj2) abc() {

}

func TestObjEq_EscapeToHeap(t *testing.T) {
	o1 := &obj1{s: "hello"}
	o2 := &obj1{s: "world"}
	assert.False(t, o1 == o2)
	var o3 obj = o1
	var o4 obj = o2
	assert.False(t, o3 == o4)

	// runtime.zerobase(SB) all zero structure object addresses
	o5 := &obj2{}
	o6 := &obj2{}
	assert.Equal(t, uintptr(0), unsafe.Sizeof(*o5))
	assert.Equal(t, uintptr(0), unsafe.Sizeof(*o6))
	t.Logf("%p", o5) // Escape to heap
	// t.Logf("%p", o6) // Escape to heap
	assert.False(t, o5 == o6)
	_o5 := &struct{}{}
	_o6 := &struct{}{}
	assert.True(t, _o5 == _o6)
	var o7 obj = o5
	var o8 obj = o6
	assert.False(t, o7 == o8)
}

func TestObjEq(t *testing.T) {
	o1 := &obj1{s: "hello"}
	o2 := &obj1{s: "world"}
	assert.False(t, o1 == o2)
	var o3 obj = o1
	var o4 obj = o2
	assert.False(t, o3 == o4)

	// runtime.zerobase(SB) all zero structure object addresses
	o5 := &obj2{}
	o6 := &obj2{}
	assert.Equal(t, uintptr(0), unsafe.Sizeof(*o5))
	assert.Equal(t, uintptr(0), unsafe.Sizeof(*o6))
	t.Logf("%p", o5)         // Escape to heap
	assert.True(t, o5 == o6) // order rearranged
	t.Logf("%p", o6)         // Escape to heap
	assert.True(t, o5 == o6)
	_o5 := &struct{}{}
	_o6 := &struct{}{}
	assert.True(t, _o5 == _o6)
	var o7 obj = o5
	var o8 obj = o6
	assert.True(t, o7 == o8)
}
