package cast

import (
	"fmt"
	"testing"
)

func TestToString(t *testing.T) {
	fmt.Println(ToString(nil))
	s := 1
	bs := []byte("aa")
	fmt.Println(ToString(s))
	fmt.Println(ToString(bs))
	fmt.Println(ToString(&bs))
	fmt.Println(ToString(&s))
	fmt.Println(ToString(&map[string]string{
		"1": "2",
	}))
	fmt.Println(ToString(1))
	fmt.Println(ToString(1.2))
	fmt.Println(ToString(-1.2))
}

func TestToInt(t *testing.T) {
	fmt.Println(ToIntI("a"))
	fmt.Println(ToIntI("1"))
	fmt.Println(ToIntI(true))
	fmt.Println(ToIntI(false))
	fmt.Println(ToIntI(1.23))
}
