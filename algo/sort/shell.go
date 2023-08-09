package sort

/*
基本思想
假设待排序的 N 个对象，首先选取一个整数作为 gap（< N），将全部对象分为 gap 个子序列，
所有距离为 gap 的对象放在同一个子序列里面，接着在每一个子序列中实施插入排序，然后逐步缩小 gap 的值，
并重复排序的过程，直到 gap == 1 且完成排序

这是一种不稳定排序，也就是不能保证相同元素之间原有的先后顺序
排序的性能取决于 gap 的数学性质（公式）

shell's original sequence: n/2, n/4, ..., 1

Knuth's increments: 1, 4, 13, ..., (3^k - 1) / 2
http://sun.aei.polsl.pl/~mciura/publikacje/shellsort.pdf

Sedgewick's increments: 1, 8, 23, 77, 281, 1073, 4193, 16577, ...,
Hibbard's increments: 1, 3, 7, 15, 31, 63, 127, 255, 511, ...
Paperno & Stasevich increments: 1, 3, 5, 9, 17, 33, 65, ...
Pratt: 1, 2, 3, 4, 6, 9, 8, 12, 18, 27, 16, 24, 36, 54, 81, ...
*/

/*
ShellSort
从 internal 不断地往 index 0 进行交换才能确保所有元素的有序性
如果是从 index 0 不断往后遍历会出现部分元素无序现象
*/
func ShellSort(arr []int) {
	n := len(arr)
	for interval := n / 2; interval > 0; interval = interval / 2 {
		for i := interval; i < n; i++ {
			temp := arr[i]
			j := i
			for ; j >= interval && arr[j-interval] > temp; j = j - interval {
				arr[j] = arr[j-interval]
			}
			arr[j] = temp
		}
	}
}
