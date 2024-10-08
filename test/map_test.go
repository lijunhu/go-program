package test

import "testing"

type Employee struct {
	Name string `json:"name"`
}

func TestMap(t *testing.T) {

	nums := []int{1, 2, 3, 4, 5}
	ans := AllCombain(nums)
	t.Log(ans)

}

func AllCombain(nums []int) (ans [][]int) {
	length := len(nums)
	for i := 0; i < length; i++ {
		ans = append(ans, nums[i:length])
	}
	return ans
}
