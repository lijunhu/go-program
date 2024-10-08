package leetcode

func mergeSortArr(nums1 []int, m int, nums2 []int, n int) []int {
	p1, p2, tail := m-1, n-1, m+n-1
	for p1 >= 0 || p2 >= 0 {
		if p1 == -1 {
			nums1[tail] = nums2[p2]
			p2--
		} else if p2 == -1 {
			nums1[tail] = nums1[p1]
			p1--
		} else if nums1[p1] > nums2[p2] {
			nums1[tail] = nums1[p1]
			p1--
		} else {
			nums1[tail] = nums2[p2]
			p2--
		}
		tail--
	}
	return nums1
}
