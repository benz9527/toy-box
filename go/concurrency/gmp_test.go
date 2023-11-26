package concurrency

import (
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestGMP(t *testing.T) {
	Init()
	t.Log(GetM(), GetMID(), GetG(), mOffset, mIDOffset)
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(5 * time.Second)
		}()
	}
	t.Log(GetM(), GetMID(), GetG(), mOffset, mIDOffset)
	wg.Wait()
	t.Log(GetM(), GetMID(), GetG(), mOffset, mIDOffset)
	t.Log(GetFieldInG[uint64]("goid"))
	SetFieldInG[uint64]("goid", 123)
	t.Log(GetFieldInG[uint64]("goid"))
	t.Log(GetFieldInP[int32]("id"))
	// Here will be a holding locks panic, if the num of P is lower than target value.
	SetFieldInP[int32]("id", 1)
	t.Log(GetFieldInP[int32]("id"))
}

func TestSuspendG(t *testing.T) {
	Init()
	targetG := make(chan uintptr)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// defer in here is never will be reachable by suspendG (i.e. park).
		targetG <- GetG()
		i := 0
		for {
			i++
			if i >= 1000 {
				i = 0
			}
		}
	}()
	go func() {
		tg := <-targetG
		time.AfterFunc(4*time.Second, func() {
			doSuspend(unsafe.Pointer(tg))
			t.Log("suspend target G")
			wg.Done()
		})
	}()
	wg.Wait()
}
