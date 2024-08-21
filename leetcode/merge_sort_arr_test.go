package leetcode

import (
	"testing"
)

func Test_merge(t *testing.T) {
	type args struct {
		nums1 []int
		m     int
		nums2 []int
		n     int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case1",
			args: args{
				nums1: []int{1, 0},
				m:     1,
				nums2: []int{1},
				n:     1,
			},
			want: []int{1, 1},
		},
		{
			name: "case2",
			args: args{
				nums1: []int{1, 1, 2, 5, 0, 0, 0, 0},
				m:     4,
				nums2: []int{2, 3, 4, 5},
				n:     4,
			},
			want: []int{1, 1, 2, 2, 3, 4, 5, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeSortArr(tt.args.nums1, tt.args.m, tt.args.nums2, tt.args.n)
			for i := range got {
				if got[i] != tt.want[i] {
					t.Error("index:", i, "got", got[i], "want", tt.want[i], "not equal")
				}
			}
		})
	}
}
