package list

import (
	"testing"
	"unsafe"
)

func TestNodeElementMemoryAlignment(t *testing.T) {
	ne1 := NewNodeElement[int](1)
	t.Logf("ne1 ptr alignment: %d", unsafe.Alignof(ne1))
	t.Logf("ne1 next ptr alignment: %d", unsafe.Alignof(ne1.(*nodeElement[int]).next))
	t.Logf("ne1 prev ptr alignment: %d", unsafe.Alignof(ne1.(*nodeElement[int]).prev))
	t.Logf("ne1 list ptr alignment: %d", unsafe.Alignof(ne1.(*nodeElement[int]).list))
	t.Logf("ne1 lock ptr alignment: %d", unsafe.Alignof(ne1.(*nodeElement[int]).lock))
	t.Logf("ne1 value alignment: %d", unsafe.Alignof(ne1.(*nodeElement[int]).value))
	t.Logf("ne1 ptr sizeof: %d", unsafe.Sizeof(ne1))
	t.Logf("ne1 next ptr sizeof: %d", unsafe.Sizeof(ne1.(*nodeElement[int]).next))
	t.Logf("ne1 prev ptr sizeof: %d", unsafe.Sizeof(ne1.(*nodeElement[int]).prev))
	t.Logf("ne1 list ptr sizeof: %d", unsafe.Sizeof(ne1.(*nodeElement[int]).list))
	t.Logf("ne1 lock ptr sizeof: %d", unsafe.Sizeof(ne1.(*nodeElement[int]).lock))
	t.Logf("ne1 value sizeof: %d", unsafe.Sizeof(ne1.(*nodeElement[int]).value))

	ne2 := nodeElement[int]{}
	t.Logf("ne2 alignment: %d", unsafe.Alignof(ne2))
	t.Logf("ne2 next ptr alignment: %d", unsafe.Alignof(ne2.next))
	t.Logf("ne2 prev ptr alignment: %d", unsafe.Alignof(ne2.prev))
	t.Logf("ne2 list ptr alignment: %d", unsafe.Alignof(ne2.list))
	t.Logf("ne2 lock ptr alignment: %d", unsafe.Alignof(ne2.lock))
	t.Logf("ne2 value alignment: %d", unsafe.Alignof(ne2.value))
	t.Logf("ne2 sizeof: %d", unsafe.Sizeof(ne2))
	t.Logf("ne2 next ptr sizeof: %d", unsafe.Sizeof(ne2.next))
	t.Logf("ne2 prev ptr sizeof: %d", unsafe.Sizeof(ne2.prev))
	t.Logf("ne2 list ptr sizeof: %d", unsafe.Sizeof(ne2.list))
	t.Logf("ne2 lock ptr sizeof: %d", unsafe.Sizeof(ne2.lock))
	t.Logf("ne2 value sizeof: %d", unsafe.Sizeof(ne2.value))

	ne3 := nodeElement[string]{}
	t.Logf("ne3 alignment: %d", unsafe.Alignof(ne3))
	t.Logf("ne3 next ptr alignment: %d", unsafe.Alignof(ne3.next))
	t.Logf("ne3 prev ptr alignment: %d", unsafe.Alignof(ne3.prev))
	t.Logf("ne3 list ptr alignment: %d", unsafe.Alignof(ne3.list))
	t.Logf("ne3 lock ptr alignment: %d", unsafe.Alignof(ne3.lock))
	t.Logf("ne3 value alignment: %d", unsafe.Alignof(ne3.value))
	t.Logf("ne3 sizeof: %d", unsafe.Sizeof(ne3))
}
