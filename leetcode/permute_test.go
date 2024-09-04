package leetcode

import "testing"

func Test_permute(t *testing.T) {
	str := "abcd"
	result := permute(str)

	for _, r := range result {
		t.Log(r)
	}

}
