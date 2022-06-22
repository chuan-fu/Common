package syncx

import (
	"fmt"
	"testing"
)

func TestAtomicInt64(t *testing.T) {
	b := NewAtomicInt64()
	fmt.Println(b.Val() == 0)

	b.Set(1)
	fmt.Println(b.Val() == 1)
	fmt.Println(b.IsVal(1))
	fmt.Println(b.IsVal(2) == false)

	b.CompareAndSwap(0, 1)
	fmt.Println(b.IsVal(1))

	b.CompareAndSwap(1, 2)
	fmt.Println(b.IsVal(2))

	b.Add(1)
	fmt.Println(b.IsVal(3))

	fmt.Println(b.Swap(4) == 3)
	fmt.Println(b.IsVal(4))

	b2 := ForAtomicInt64(6)
	fmt.Println(b2.IsVal(6))
}
