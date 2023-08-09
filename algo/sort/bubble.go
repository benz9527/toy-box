package sort

func BubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				temp := arr[i]
				arr[i] = arr[j]
				arr[j] = temp
			}
		}
	}
}

func BubbleSortAcc(arr []int) {
	n := len(arr)
	doSort := false
	for i := 0; i < n; i++ {
		doSort = false
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				temp := arr[i]
				arr[i] = arr[j]
				arr[j] = temp
				doSort = true
			}
		}

		if !doSort {
			continue
		}
	}
}
