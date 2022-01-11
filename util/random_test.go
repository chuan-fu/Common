package util

import (
	"fmt"
	"testing"
)

func TestTT(*testing.T) {
	r := NewRandomParam(AllChar)
	for i := 0; i < 100; i++ {
		fmt.Println(r.GenRandomKey(6))
	}
}
