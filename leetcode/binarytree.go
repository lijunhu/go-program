package leetcode

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func preorderTraversal(root *TreeNode) []int {
	var result []int
	var traverse func(node *TreeNode)

	traverse = func(node *TreeNode) {
		if node == nil {
			return
		}
		result = append(result, node.Val)
		traverse(node.Left)
		traverse(node.Right)
	}

	traverse(root)
	return result
}

func preorderTraversalStack(root *TreeNode) []int {
	var result []int
	if root == nil {
		return result
	}

	stack := []*TreeNode{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, node.Val)

		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}

	return result
}

func inorderTraversal(root *TreeNode) []int {
	var result []int
	var traverse func(node *TreeNode)

	traverse = func(node *TreeNode) {
		if node == nil {
			return
		}
		traverse(node.Left)
		result = append(result, node.Val)
		traverse(node.Right)
	}

	traverse(root)
	return result
}

func inorderTraversalStack(root *TreeNode) []int {
	var result []int
	var stack []*TreeNode
	current := root

	for current != nil || len(stack) > 0 {
		for current != nil {
			stack = append(stack, current)
			current = current.Left
		}
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, current.Val)
		current = current.Right
	}

	return result
}

func postorderTraversal(root *TreeNode) []int {
	var result []int
	var traverse func(node *TreeNode)

	traverse = func(node *TreeNode) {
		if node == nil {
			return
		}
		traverse(node.Left)
		traverse(node.Right)
		result = append(result, node.Val)
	}

	traverse(root)
	return result
}

func postorderTraversalStack(root *TreeNode) []int {
	var result []int
	if root == nil {
		return result
	}

	stack := []*TreeNode{root}
	var prev *TreeNode
	for len(stack) > 0 {
		current := stack[len(stack)-1]

		if prev == nil || prev.Left == current || prev.Right == current {
			if current.Left != nil {
				stack = append(stack, current.Left)
			} else if current.Right != nil {
				stack = append(stack, current.Right)
			}
		} else if current.Left == prev {
			if current.Right != nil {
				stack = append(stack, current.Right)
			}
		} else {
			result = append(result, current.Val)
			stack = stack[:len(stack)-1]
		}

		prev = current
	}

	return result
}

func levelOrder(root *TreeNode) [][]int {
	var result [][]int
	if root == nil {
		return result
	}

	queue := []*TreeNode{root}
	for len(queue) > 0 {
		levelSize := len(queue)
		var level []int
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, level)
	}

	return result
}
