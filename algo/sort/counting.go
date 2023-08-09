package sort

/*
CountingSort
It sorts the elements of an array by counting the number of occurrences of each unique element in the array.
The count is stored in an auxiliary array and the sorting si done by mapping the count as an index of the auxiliary array.

1.Find out the maximum element (let ie be maxElement) from the given array.
2.Initialize an array of length as maxElement value + 1 with all elements 0. This array is used for storing the count of
the elements in the array.
3.Store the count of each element at their respective index in count array
4.Store cumulative sum of the elements of the count array. It helps in placing the elements into the correct index of
the sorted array.
5.Find the index of each element of the original array in the count array. This gives the cumulative count. Place the
element at the index calculated.
6.After placing each element at its correct position, decrease its count by one.
*/
func CountingSort(arr []int) {
	n := len(arr)
	maxElem := arr[0]
	for i := 1; i < n; i++ {
		if arr[i] > maxElem {
			maxElem = arr[i]
		}
	}

	maxLen := maxElem + 1
	auxCountingArr := make([]int, maxLen)

	for i := 0; i < n; i++ {
		auxCountingArr[arr[i]] = auxCountingArr[arr[i]] + 1
	}

	originIdx := 0
	for i := 0; i < maxLen; i++ {
		for auxCountingArr[i] > 0 {
			arr[originIdx] = i
			auxCountingArr[i] = auxCountingArr[i] - 1
			originIdx = originIdx + 1
		}
	}
}
