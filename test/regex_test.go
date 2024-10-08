package test

import (
	"regexp"
	"testing"
)

const methodPattern = "^[a-zA-Z0-9\\.]*$"

func TestRegex(t *testing.T) {

	str := "a_2.4"

	r, err := regexp.Compile(methodPattern)
	if err != nil {
		t.Log(err)
	}
	t.Log(r.MatchString("a_b.c"))

	rep, err := regexp.Compile(`^[(a-zA-Z0-9).\\-_]+[a-zA-Z0-9]$`)
	if err != nil {
		t.Log(err)
	}
	t.Log(rep.MatchString(str))
}

type Node struct {
	val  int
	next *Node
}

func mergeListNode(node1, node2 *Node) (ret *Node) {
	if node1 == nil {
		return node2
	}
	if node2 == nil {
		return node1
	}
	curr := node1
	tail := node2
	for curr != nil {
		if curr.val < tail.val {

		}
		curr = curr.next
	}
	return
}

func reverseListNode(node *Node) (ret *Node) {
	if node == nil {
		return nil
	}
	curr := node
	tail := new(Node)
	for curr != nil {
		next := curr.next
		curr.next = tail
		tail = curr
		curr = next
	}
	return tail
}

func Test_ReverseListNode(t *testing.T) {

	//var m =make(map[string]string)

	n1 := &Node{1, nil}
	n2 := &Node{2, n1}
	n3 := &Node{3, n2}
	n4 := &Node{4, n3}
	n5 := &Node{5, n4}

	n := reverseListNode(n5)
	t.Log(n.val)
}
