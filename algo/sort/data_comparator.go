package sort

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type sortFunc func([]int)

func GenerateRandomArr(maxSize int, maxNum int) []int {
	if maxSize <= 0 {
		panic("incorrect max size for array")
	}

	var arr []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < maxSize; i++ {
		arr = append(arr, rand.Intn(maxNum))
	}

	return arr
}

func IsEquals(standardArr []int, customizedArr []int) bool {
	if (standardArr != nil && customizedArr == nil) ||
		(standardArr == nil && customizedArr != nil) {
		return false
	}

	if len(standardArr) != len(customizedArr) {
		return false
	}

	for i := 0; i < len(standardArr); i++ {
		if standardArr[i] != customizedArr[i] {
			return false
		}
	}
	return true
}

func CopyArr(originArr []int, maxSize int) []int {
	if originArr == nil {
		return nil
	}
	dstArr := make([]int, maxSize)
	copy(dstArr, originArr)
	return dstArr
}

func PrintArr(arr []int) {
	if arr == nil {
		return
	}

	fmt.Println(arr)
}

func DataComparator(sf sortFunc, benchmarkNum int, maxArrSize int, maxNum int, isAsc bool) {
	equals := true
	for i := 0; i < benchmarkNum; i++ {
		originArr := GenerateRandomArr(maxArrSize, maxNum)
		arr1 := CopyArr(originArr, maxArrSize)
		arr2 := CopyArr(originArr, maxArrSize)

		if isAsc {
			sort.Slice(arr1, func(i, j int) bool {
				return arr1[i] < arr1[j]
			})
		} else {
			sort.Slice(arr1, func(i, j int) bool {
				return arr1[i] > arr1[j]
			})
		}
		sf(arr2)

		if !IsEquals(arr1, arr2) {
			PrintArr(originArr)
			PrintArr(arr1)
			PrintArr(arr2)
			equals = false
			break
		}
	}

	if equals {
		fmt.Println("Success correct")
	} else {
		fmt.Println("Error")
	}
}
