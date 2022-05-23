package util

import (
	"fmt"
	"testing"
)

func TestBytesToString(t *testing.T) {
	fmt.Println(BytesToString(nil))
	fmt.Println(BytesToString([]byte("")))
	fmt.Println(BytesToString([]byte("1")))
}

func TestInt64Join(t *testing.T) {
	fmt.Println(Int64Join(nil, ",,"))
	fmt.Println(Int64Join([]int64{}, ","))
	fmt.Println(Int64Join([]int64{1, 2, 3}, ","))
	fmt.Println(Int64Join([]int64{1, 2, 3, 1}, ","))
}
