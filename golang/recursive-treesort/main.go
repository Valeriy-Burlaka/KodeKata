package main

import (
	"fmt"
)

type tree struct {
	value       int
	left, right *tree
}

func add(t *tree, value int) *tree {
	if t == nil {
		return &tree{value: value}
	}

	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}

	return t
}

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}

	return values
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

func main() {
	s := []int{7, 4, 4, 2, 9, 2, 3, 4, 6, 7}
	Sort(s)
	fmt.Println(s)
	// [2 2 3 4 4 4 6 7 7 9]
}
