package util

import (
	"fmt"
	"testing"
)

func TestToString(t *testing.T) {
	fmt.Println(ToString(map[string]string{
		"1": "2",
	}))
	fmt.Println(ToString(1))
	fmt.Println(ToString(1.2))
	fmt.Println(ToString(-1.2))
}
