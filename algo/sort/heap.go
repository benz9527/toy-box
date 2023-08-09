package sort

/*
1.since the tree satisfies Max-Heap property, then the largest item is stored at
the root node.
2.swap: remove the root element and put at the end of the array (nth position)
Put the last item of the tree (heap) at the vacant place.
3.remove: reduce the size of the heap by 1
4.heapify: heapify the root element again so that we have the highest element at root
5.the process is repeated until all the items of the list are sorted.
*/

func heapify(arr []int, n int, i int) {
	// find largest among root, left child and right child
	// such as index 0, its left child index is 1 and right child index is 2
	largest := i
	l := 2*i + 1
	r := 2 * (i + 1)

	if l < n && arr[l] > arr[largest] {
		largest = l
	}

	if r < n && arr[r] > arr[largest] {
		largest = r
	}

	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]

		heapify(arr, n, largest)
	}
}

/*
总是从第一层非叶子节点开始遍历
*/
func HeapSort(arr []int) {
	n := len(arr)
	// build max heap
	// to build a max-heap from any tree, we can thus start heapifying each sub-tree
	// from the bottom up and end up with a max-heap after the function is applied to
	// all the elements including the root element.
	// in the case of a complete tree. the first index of a non-leaf node is given by
	// n / 2 - 1.
	// all other nodes after that are leaf-nodes and thus don't need to be heapified.
	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}

	for i := n - 1; i >= 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]

		// heapify root element
		heapify(arr, i, 0)
	}
}
