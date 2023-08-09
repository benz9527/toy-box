package sort

import (
	"math/rand"
	"time"
)

func swap(arr []int, i int, j int) {
	temp := arr[i]
	arr[i] = arr[j]
	arr[j] = temp
}

func equalPartition(arr []int, l int, r int) (int, int) {
	less := l - 1
	more := r
	pivot := arr[r]

	for more > l {
		if arr[l] < pivot {
			less++
			if l != less {
				swap(arr, l, less)
			}
			l++
		} else if arr[l] > pivot {
			more--
			if l != more {
				swap(arr, l, more)
			}
		} else {
			l++
		}
	}

	if more != r {
		swap(arr, more, r)
	}
	return less + 1, more - 1
}

func quickSort(arr []int, l int, r int) {
	if r > l {
		swap(arr, r, l+rand.Intn(r-l+1))
		loc1, loc2 := equalPartition(arr, l, r)
		quickSort(arr, l, loc1-1)
		quickSort(arr, loc2+1, r)
	}
}

func QuickSort(arr []int) {
	n := len(arr)
	if arr == nil || n < 2 {
		return
	}
	rand.Seed(time.Now().UnixNano())
	quickSort(arr, 0, n-1)
}

/*  another implementation  */
func partition(arr []int, low int, high int) int {
	// choose the rightmost element as pivot
	pivot := arr[high]

	// pointer for greater element
	i := low - 1

	// traverse through all elements and compare each element with pivot
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			// if element smaller than pivot is found swap it with the greater element pointed by i
			i++

			// swapping element at i with element at j
			temp := arr[i]
			arr[i] = arr[j]
			arr[j] = temp
		}
	}

	// swap the pivot element with the greater element specified by i
	temp := arr[i+1]
	arr[i+1] = arr[high]
	arr[high] = temp

	// return the position from where partition is done
	return i + 1
}

func quickSort2(arr []int, low int, high int) {
	if low < high {
		// find pivot element such that
		// elements smaller than pivot are on the left
		// elements greater that pivot are on the right
		pi := partition(arr, low, high)

		quickSort2(arr, low, pi-1)
		quickSort2(arr, pi+1, high)
	}
}

func QuickSort2(arr []int) {
	n := len(arr)
	if arr == nil || n < 2 {
		return
	}
	quickSort2(arr, 0, n-1)
}
