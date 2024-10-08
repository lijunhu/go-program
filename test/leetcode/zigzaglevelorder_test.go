package leetcode

import (
	"fmt"
	"testing"
)

func TestZigzagLevelOrder(t *testing.T) {
	root := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 4,
			},
			Right: &TreeNode{
				Val: 5,
			},
		},
		Right: &TreeNode{
			Val: 3,
			Left: &TreeNode{
				Val: 6,
			},
			Right: &TreeNode{
				Val: 7,
			},
		},
	}

	result := ZigzagLevelOrder(root)
	fmt.Println(result) // Output: [[1], [3, 2], [4, 5, 6, 7]]
}
