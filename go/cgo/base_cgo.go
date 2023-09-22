package cgo

/*
#include <stdio.h>

void printi(int v) {
    printf("printi: %d\n", v);
}
*/
import "C"

func printByC(num int) {
	C.printi(C.int(num))
	return
}
