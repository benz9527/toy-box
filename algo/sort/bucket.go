package sort

/*
The process of bucket sort can be understood as scatter-gather approach.
Here, elements are first scattered into buckets then the elements in each buckets are sorted.
Finally, the elements are gathered in order.

1.创建固定数量的 buckets
2.bucket size 最好通过归一化的处理后平均划分（也不是绝对均分），大于 0 小于 1 的 float 可以通过扩大相同倍数（> 1 整数）
来进行 size 划分
*/

func BucketSort(arr []int) {
	n := len(arr)
	BUCKET_SIZE := 10
	buckets := make([][]int, BUCKET_SIZE)
	max := getMax(arr, n)
	var interval int = max / BUCKET_SIZE
	if max%BUCKET_SIZE > 0 {
		interval++
	}

	for i := 0; i < n; i++ {
		bucketIdx := arr[i]/interval - 1
		if arr[i]%interval > 0 {
			bucketIdx++
		} else if bucketIdx == -1 {
			bucketIdx = 0
		}

		buckets[bucketIdx] = append(buckets[bucketIdx], arr[i])
	}

	origIdx := 0
	for i := 0; i < BUCKET_SIZE && origIdx < n; i++ {
		blen := len(buckets[i])
		if blen > 0 {
			InsertSortSimplify(buckets[i])

			for j := 0; j < blen; j++ {
				arr[origIdx] = buckets[i][j]
				origIdx++
			}
		}
	}
}
