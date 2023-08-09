package sort

/*
https://www.programiz.com/dsa/merge-sort
Based on the principle of Divide and Conquer Algorithm
*/

func divideAndConquer(arr []int, leftIdx int, rightIdx int) {
	if leftIdx >= rightIdx {
		return
	}

	middleIdx := (leftIdx + rightIdx) / 2

	divideAndConquer(arr, leftIdx, middleIdx)
	divideAndConquer(arr, middleIdx+1, rightIdx)

	merge(arr, leftIdx, middleIdx, rightIdx)
}

func merge(arr []int, subl int, subm int, subr int) {
	n1 := subm - subl + 1
	n2 := subr - subm

	arrL := []int{}
	arrR := []int{}

	for i := 0; i < n1; i++ {
		arrL = append(arrL, arr[subl+i])
	}

	for i := 0; i < n2; i++ {
		arrR = append(arrR, arr[subm+1+i])
	}

	i, j, k := 0, 0, subl

	for i < n1 && j < n2 {
		if arrL[i] <= arrR[j] {
			arr[k] = arrL[i]
			i++
		} else {
			arr[k] = arrR[j]
			j++
		}
		k++
	}

	for i < n1 {
		arr[k] = arrL[i]
		i++
		k++
	}

	for j < n2 {
		arr[k] = arrR[j]
		j++
		k++
	}
}

func MergeSort(arr []int) {
	n := len(arr)
	divideAndConquer(arr, 0, n-1)
}
