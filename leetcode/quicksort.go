package leetcode

func quickSort(nums []int, low, high int) {
	if len(nums) <= 1 {
		return
	}
	if low >= high {
		return
	}
	mid := partition(nums, low, high)
	quickSort(nums, low, mid-1)
	quickSort(nums, mid+1, high)
}

// partition 将数组分区并返回基准点的索引
func partition(arr []int, low, high int) int {
	pivot := arr[high] // 选择最后一个元素作为基准
	i := low - 1       // 初始化较小元素的索引

	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i] // 交换元素
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1] // 将基准元素移到正确位置
	return i + 1
}
