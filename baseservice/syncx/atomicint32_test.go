package syncx

import (
	"fmt"
	"sync"
	"testing"
)

func TestAtomicInt32(t *testing.T) {
	b := NewAtomicInt32()
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

	b2 := ForAtomicInt32(6)
	fmt.Println(b2.IsVal(6))
}

func TestAtomicInt32Add(t *testing.T) {
	b := NewAtomicInt32()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				b.AddAtomic(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(b.Val())
}
