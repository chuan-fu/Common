package util

import "testing"

func TestName(t *testing.T) {
	a := "aaa"
	defer DeferFunc()
	panic(&a)
}
