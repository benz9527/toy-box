//go:build !gccgo

// https://wiki.archlinux.org/title/Talk:Go
//The gccgo compiler (from gcc-go) is usually slower to compile but generally has more optimizations than the gc compiler (from go).
//
//The Go language has always been defined by a spec, not an implementation. The Go team has written two different compilers that implement that spec: gc and gccgo.
//Gc is the original compiler, and the go tool uses it by default.
//Gccgo is a different implementation with a different focus
//Compared to gc, gccgo is slower to compile code but supports more powerful optimizations, so a CPU-bound program built by gccgo will usually run faster.
//The gc compiler supports only the most popular processors: x86 (32-bit and 64-bit) and ARM.
//Gccgo, however, supports all the processors that GCC supports.
//Not all those processors have been thoroughly tested for gccgo, but many have, including x86 (32-bit and 64-bit), SPARC, MIPS, PowerPC and even Alpha.
//Gccgo has also been tested on operating systems that the gc compiler does not support, notably Solaris.
//— Gccgo in GCC 4.7.1 (July 2012, slightly reformatted for clarity)
//-- The Puzzlemaker (talk) 04:46, 8 March 2020 (UTC)
//
//Recent 2019 benchmarks appear to show that gc has improved more quickly than gccgo has (probably because it's received more attention), so even though gccgo has made significant improvements, gc seems like the better choice in most circumstances: https://meltware.com/2019/01/16/gccgo-benchmarks-2019.html Skyfaller (talk) 12:52, 11 April 2020 (UTC)

package concurrency

import (
	"unsafe"

	"github.com/modern-go/reflect2"
	_ "github.com/v2pro/plz/gls"
)

// Refers: https://github.com/MeteorsLiu/getm/blob/main/getg.go

var (
	gType     reflect2.StructType
	mType     reflect2.StructType
	pType     reflect2.StructType
	mOffset   uintptr
	pOffset   uintptr
	mIDOffset uintptr
)

func Init() {
	gType = reflect2.TypeByName("runtime.g").(reflect2.StructType)
	if gType == nil {
		panic("failed to get runtime.g type")
	}
	mType = reflect2.TypeByName("runtime.m").(reflect2.StructType)
	if mType == nil {
		panic("failed to get runtime.m type")
	}
	pType = reflect2.TypeByName("runtime.p").(reflect2.StructType)
	if pType == nil {
		panic("failed to get runtime.p type")
	}
	mOffset = gType.FieldByName("m").Offset()
	pOffset = mType.FieldByName("p").Offset()
	mIDOffset = mType.FieldByName("id").Offset()
}

//go:linkname GetG github.com/v2pro/plz/gls.getg
func GetG() uintptr

func GetM() uintptr {
	g := GetG()
	m := (*uintptr)(unsafe.Pointer(g + mOffset))
	return *m
}

func GetP() uintptr {
	m := GetM()
	p := (*uintptr)(unsafe.Pointer(m + pOffset))
	return *p
}

func GetMID() int64 {
	m := GetM()
	mid := (*int64)(unsafe.Pointer(m + mIDOffset))
	return *mid
}

// Experimental APIs, high risk and unsafe.
// Most likely occur panic.

func GetFieldInG[T any](fieldName string) T {
	customizedOffset := gType.FieldByName(fieldName).Offset()
	return *(*T)(unsafe.Pointer(GetG() + customizedOffset))
}

func GetFieldInM[T any](fieldName string) T {
	customizedOffset := mType.FieldByName(fieldName).Offset()
	return *(*T)(unsafe.Pointer(GetM() + customizedOffset))
}

func GetFieldInP[T any](fieldName string) T {
	customizedOffset := pType.FieldByName(fieldName).Offset()
	return *(*T)(unsafe.Pointer(GetP() + customizedOffset))
}

func SetFieldInG[T any](fieldName string, value T) {
	customizedOffset := gType.FieldByName(fieldName).Offset()
	*(*T)(unsafe.Pointer(GetG() + customizedOffset)) = value
}

func SetFieldInM[T any](fieldName string, value T) {
	customizedOffset := mType.FieldByName(fieldName).Offset()
	*(*T)(unsafe.Pointer(GetM() + customizedOffset)) = value
}

func SetFieldInP[T any](fieldName string, value T) {
	customizedOffset := pType.FieldByName(fieldName).Offset()
	*(*T)(unsafe.Pointer(GetP() + customizedOffset)) = value
}

// runtime.gopark
// mcall(park_m)
// func mcall(f func(*g)) {
// 	  g := getg()
// 	  f(g)
//	}
// go1.21.3/src/runtime/proc.go:3721
// https://cs.opensource.google/go/go/+/refs/heads/master:src/runtime/preempt.go
// Yield the M to other goroutines.
//
// 1. time.Sleep
// gopark(resetForSleep, unsafe.Pointer(t), waitReasonSleep, traceBlockSleep, 1)
// 2. channel read or write wait block
// gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceBlockChanSend, 2)
// 3. sync.WaitGroup
// 4. sync.Mutex, sync.RWMutex, sync.Cond
// goparkunlock
// gopark(parkunlock_c, unsafe.Pointer(lock), reason, traceReason, traceskip)
// 5. network io read or write block
// netpoll
// gopark(netpollblockcommit, unsafe.Pointer(gpp), waitReasonIOWait, traceBlockNet, 5)

//go:linkname suspendG runtime.suspendG
func suspendG(unsafe.Pointer)

//go:linkname systemstack runtime.systemstack
func systemstack(fn func())

//go:linkname casGToWaiting runtime.casGToWaiting
func casGToWaiting(gp unsafe.Pointer, old uint32, reason uint8)

//go:linkname casgstatus runtime.casgstatus
func casgstatus(gp unsafe.Pointer, oldval, newval uint32)

func doSuspend(g unsafe.Pointer) {
	systemstack(func() {
		_g := unsafe.Pointer(GetFieldInM[uintptr]("curg"))
		// 2: _Grunning to 4: _Gwaiting
		// reason 7: waitReasonGarbageCollectionScan
		casGToWaiting(_g, 2, 7)
		suspendG(g)
		// 4: _Gwaiting to 2: _Grunning
		casgstatus(_g, 4, 2)
	})
}
