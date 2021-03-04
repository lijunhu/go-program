package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestSliceEqual(t *testing.T) {
	nameEnv := "sz.qa_officeidc_hd1"
	fmt.Println(nameEnv[strings.Index(nameEnv, ".")+1:strings.Index(nameEnv, "_")])

	a := []Rule{{Value: "1", KeyName: "2", KeyType: "3", Type: "4"}}
	b := []Rule{{Value: "1", KeyName: "2", KeyType: "3", Type: "4"}}

	fmt.Println(ObjectSliceEqual(a, b))
}

func ObjectSliceEqual(a, b []Rule) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
