package syncx

import (
	"fmt"
	"testing"
)

func TestAtomicbool(t *testing.T) {
	b := NewAtomicBool()
	fmt.Println(b.True() == false)

	b.Set(true)
	fmt.Println(b.True() == true)

	b.CompareAndSwap(false, true)
	fmt.Println(b.True() == true)

	b.CompareAndSwap(true, false)
	fmt.Println(b.True() == false)

	b2 := ForAtomicBool(true)
	fmt.Println(b2.True() == true)

	b.Set(false)
	fmt.Println(b.True() == false)
}
