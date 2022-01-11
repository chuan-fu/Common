package util

import (
	"fmt"
	"testing"
)

func TestTT(*testing.T) {
	r := NewRandomParam(0)
	for i := 0; i < 100; i++ {
		fmt.Println(r.GenRandomKey(6))
	}
}
