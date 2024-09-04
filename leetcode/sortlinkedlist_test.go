package leetcode

import (
	"testing"
)

func Test_mergeSort(t *testing.T) {
	// 创建链表: 4 -> 2 -> 1 -> 3
	head := &ListNode{Val: 4}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 1}
	head.Next.Next.Next = &ListNode{Val: 3}

	t.Log("排序前:")
	printList(head)

	// 对链表进行归并排序
	sortedHead := mergeSort(head)

	t.Log("排序后:")
	printList(sortedHead)
}
