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
