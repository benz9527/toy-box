package sort

/*
Sorts the elements by first grouping the individual digits of the same place value.
Then, sort the elements according to their increasing / decreasing order.

1. Find the largest element in the array, i.e. max. Let X be the number of digits in max. X is calculated because
we have to go through all the significant places of all elements.
In this array [121, 1, 788], we have the largest number 788. It has 3 digits. Therefore, the loop should go up to hundreds
place (3 times).
2.Go through each significant place one by one.
*/

func getMax(arr []int, n int) (max int) {
	max = arr[0]
	for i := 1; i < n; i++ {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return max
}

func countingSort(arr []int, size int, place int) {
	output := make([]int, size+1)
	// 基数桶
	count := make([]int, 10)

	// calculate count of elements 计算数字的位数并存储
	// 0-9 的索引代表了所有数字拆分后的各个位的具体数值， 由于默认使得 0-9 的 count array 是有序的，而且是稳定有序的
	// count array 每在一个地方 + 1 表示就多一个元素
	for i := 0; i < size; i++ {
		// count[(arr[i] / place) % 10] = count[(arr[i] / place) % 10] + 1
		count[(arr[i]/place)%10]++
	}

	// calculate cumulative cont
	// 计算基数桶的边界索引，count[i] 的值为第 i 个桶的右边界索引 +1
	// 这里不断累加也就是为了确定 count[i] 中元素的索引分布范围（可能包含多个相同的元素）
	// count[i] 包括左侧元素的范围和本身元素的范围
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	// place the elements in sorted order
	// 从右往左扫描，保证排序的稳定性
	for i := size - 1; i >= 0; i-- {
		// 元素个数转为索引必然要 - 1
		output[count[(arr[i]/place)%10]-1] = arr[i]
		// 每当一个元素找到直接的位置后 count[i] 就会减一
		count[(arr[i]/place)%10]--
	}

	for i := 0; i < size; i++ {
		arr[i] = output[i]
	}
}

func RadixSort(arr []int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	max := getMax(arr, n)

	for place := 1; max/place > 0; place *= 10 {
		countingSort(arr, n, place)
	}
}
