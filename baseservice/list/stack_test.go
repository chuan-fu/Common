package list

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack()
	s.Push(11)
	s.Push(22)
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
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
