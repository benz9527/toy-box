package main

import "fmt"

// [GOSSAFUNC=main] go build -gcflags=-S main.go
// Please don't run `go tool compile -N -l -S main.go` directly.
// Run `go build -x -work main.go 1> transcript.txt 2>&1`

// References:
// https://github.com/golang/go/issues/58629

func f1() {
	fmt.Println("Ben Zheng")
}

func main() {
	f1()
}
