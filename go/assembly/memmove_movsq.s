#include "textflag.h"

// func copy_movsq(to, from unsafe.Pointer, n uintptr) (left, copied uintptr)
TEXT ·copy_movsq(SB),NOSPLIT,$0
    MOVQ to+0(FP), DI
    MOVQ from+8(FP), SI
    MOVQ n+16(FP), BX
    MOVQ BX, AX
    MOVQ BX, CX
    SHRQ $3, CX
    ANDQ $7, BX
    REP; MOVSQ
    SUBQ BX, AX
    MOVQ BX, left+24(FP)
    MOVQ AX, copied+32(FP)
    RET
