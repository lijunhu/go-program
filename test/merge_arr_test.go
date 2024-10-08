package test

import (
	"fmt"
	"testing"
)

func merge(arr1, arr2 []int)(ret []int) {
	m := len(arr1)
	n := len(arr2)
	max := m + n-1
	ret = make([]int,max+1)
	i := m - 1
	j := n -1
	for i >= 0 && j >= 0 {
		if arr1[i] > arr2[j] {
			ret[max] = arr1[i]

			// i 向左移动
			i--
		} else {
			ret[max] = arr2[j]

			// j 向左移动
			j--
		}

		// max 向左移动
		max --

	}

	for j >= 0 {
		ret[max] = arr2[j]
		max--
		j--
	}

	for i >= 0 {
		ret[max] = arr1[i]
		max--
		i--
	}
	return ret
}

func TestMerge(t *testing.T)  {

	fmt.Println(merge([]int{1,2,3,4,5},[]int{6,7,8,9}))

}