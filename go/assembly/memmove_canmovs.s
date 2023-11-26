#include "textflag.h"

// func copy_movs(to, from unsafe.Pointer, n uintptr) (ok bool)
TEXT ·can_movs(SB),NOSPLIT,$0
    MOVQ to+0(FP), DI
    MOVQ from+8(FP), SI
    MOVQ n+16(FP), CX
    // DI <= SI
    CMPQ SI, DI
    JGE true
    // SI + n <= DI
    ADDQ CX, SI
    CMPQ SI, DI
    JLS true
    MOVB $0, ok+24(FP)
    RET
true:
    MOVB $1, ok+24(FP)
    RET
