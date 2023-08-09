package sort

func InsertSort(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		edgeIdx := i - 1
		if arr[edgeIdx] <= arr[i] {
			continue
		} else {
			for j := i; j >= 1 && arr[j] < arr[j-1]; j-- {
				temp := arr[j]
				arr[j] = arr[j-1]
				arr[j-1] = temp
			}
		}
	}
}

func InsertSortSimplify(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		value := arr[i]
		edgeIdx := i
		for edgeIdx >= 1 && arr[edgeIdx-1] > value {
			arr[edgeIdx] = arr[edgeIdx-1]
			edgeIdx--
		}
		arr[edgeIdx] = value
	}
}
