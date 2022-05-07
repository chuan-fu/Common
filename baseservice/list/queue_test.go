package list

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	s := NewQueue()
	s.Push(11)
	s.Push(22)
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Head())
	fmt.Println(s.Tail())
	s.Push(33)
	s.Push(44)
	s.Push(55)
	fmt.Println(s.Len())
	fmt.Println(s.Head())
	fmt.Println(s.Tail())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
}
