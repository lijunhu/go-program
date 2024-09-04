package leetcode

import (
	"testing"
)

func Test_longestCommonSubsequence(t *testing.T) {
	t1 := "abcdefghijklmnopqrstuvwxyz"

	t2 := "123123123123abcdefghijklmnopqrstuvwxyzdasdasdasda"

	subLength := longestCommonSubsequence(t1, t2)
	t.Log(subLength)

}
