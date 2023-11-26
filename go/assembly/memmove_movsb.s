#include "textflag.h"

// https://stackoverflow.com/questions/27804852/assembly-rep-movs-mechanism
// MOVS instruction:
// Copied data from ES:(SI/EDI/RDI) to DS:(DI/ESI/RSI) and increases/decreases SI and DI
// https://c9x.me/x86/html/file_module_x86_id_279.html
// It cannot use other registers than SI/DI.
// Slowest.
// 2013 intel revisit REP MOVS again and implemented CPUID ERMSB (Enhanced REP MOVSB and STOSB Bit operation) bit.
// It is faster than REP MOVS.
// http://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-optimization-manual.pdf
// Please note that ERMSB produces best results for REP MOVSB, not REP MOVSD (MOVSQ).
// MOVSD on 32-bit
// MOVSQ on 64-bit

// func copy_movsb(to, from unsafe.Pointer, n uintptr)
TEXT ·copy_movsb(SB),NOSPLIT,$0
    MOVQ to+0(FP), DI
    MOVQ from+8(FP), SI
    MOVQ n+16(FP), CX
    REP; MOVSB // REP: Repeats the RCX times;
               // MOVSB: Copies data from RSI to RDI and increments/decrements RSI and RDI
    RET
