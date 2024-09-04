package leetcode

import (
	"fmt"
	"testing"
)

func Test_quickSort(t *testing.T) {
	arr := []int{29, 10, 14, 37, 13}
	fmt.Println("Original array:", arr)

	quickSort(arr, 0, len(arr)-1)

	fmt.Println("Sorted array:", arr)
}
