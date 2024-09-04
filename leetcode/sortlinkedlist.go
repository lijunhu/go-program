package leetcode

import (
	"fmt"
)

// 定义链表节点
type ListNode struct {
	Val  int
	Next *ListNode
}

// 归并排序的主函数
func mergeSort(head *ListNode) *ListNode {
	// 基础情况: 如果链表为空或只有一个节点，直接返回
	if head == nil || head.Next == nil {
		return head
	}

	// 使用快慢指针找到链表中点
	mid := getMid(head)
	left := head
	right := mid.Next
	mid.Next = nil

	// 递归排序左右两部分
	left = mergeSort(left)
	right = mergeSort(right)

	// 合并两个有序链表
	return merge(left, right)
}

// 使用快慢指针找到链表的中间节点
func getMid(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

// 合并两个有序链表
func merge(left, right *ListNode) *ListNode {
	dummy := &ListNode{}
	curr := dummy

	for left != nil && right != nil {
		if left.Val < right.Val {
			curr.Next = left
			left = left.Next
		} else {
			curr.Next = right
			right = right.Next
		}
		curr = curr.Next
	}

	if left != nil {
		curr.Next = left
	} else {
		curr.Next = right
	}

	return dummy.Next
}

// 打印链表
func printList(head *ListNode) {
	for head != nil {
		fmt.Printf("%d -> ", head.Val)
		head = head.Next
	}
	fmt.Println("nil")
}

// 示例使用
func main() {
	// 创建链表: 4 -> 2 -> 1 -> 3
	head := &ListNode{Val: 4}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 1}
	head.Next.Next.Next = &ListNode{Val: 3}

	fmt.Println("排序前:")
	printList(head)

	// 对链表进行归并排序
	sortedHead := mergeSort(head)

	fmt.Println("排序后:")
	printList(sortedHead)
}
