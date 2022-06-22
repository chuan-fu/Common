package util

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
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

func TestFFF(t *testing.T) {
	fmt.Println(f())
}

func f() (err error) {
	v, err := g()
	if err != nil {
		return
	}
	_ = v
	return
}

func g() (bool, error) {
	return true, errors.New("A")
}
