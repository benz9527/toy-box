package sort

import (
	"fmt"
	"testing"
)

func TestGenerateRandomArr(t *testing.T) {
	fmt.Println(GenerateRandomArr(10, 100))
}

func TestCopyArr(t *testing.T) {
	origArr := GenerateRandomArr(10, 100)
	fmt.Printf("orig arr addr: %p\n", &origArr)

	arr1 := CopyArr(origArr, len(origArr))
	fmt.Printf("arr1 addr: %p\n", &arr1)

	arr2 := CopyArr(origArr, len(origArr))
	fmt.Printf("arr2 addr: %p\n", &arr2)

}

func TestBubbleSort(t *testing.T) {
	arr := []int{6, 4, 8, 9, 0, 1, 3, 5, 2, 7}
	BubbleSort(arr)
	fmt.Println(arr)
}

func TestBubbleSortBenchmark(t *testing.T) {
	LogarithmicDetector(BubbleSort, 1000, 20, 100, true)
	LogarithmicDetector(BubbleSortAcc, 1000, 20, 100, true)
}

func TestBubbleSortAcc(t *testing.T) {
	arr := []int{6, 4, 8, 9, 0, 1, 3, 5, 2, 7}
	BubbleSortAcc(arr)
	fmt.Println(arr)
}
