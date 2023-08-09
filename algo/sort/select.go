package sort

func SelectSort(arr []int) {
	n := len(arr)
	for i := 0; i < n; i++ {
		minIdx := i
		tempVal := arr[i]
		for j := i + 1; j < n; j++ {
			if arr[j] < tempVal {
				minIdx = j
				tempVal = arr[j]
			}
		}

		arr[i], arr[minIdx] = arr[minIdx], arr[i]
	}
}
