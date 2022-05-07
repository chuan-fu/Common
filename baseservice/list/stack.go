package list

// 栈
type Stack struct {
	link *LinkList
}

func NewStack() *Stack {
	return &Stack{NewLinkList()}
}

func (s *Stack) Push(v interface{}) {
	s.link.PushFront(v)
}

func (s *Stack) Pop() interface{} {
	return s.link.PopFront()
}

func (s *Stack) Head() interface{} {
	return s.link.Head()
}

func (s *Stack) Tail() interface{} {
	return s.link.Tail()
}

func (s *Stack) Len() int {
	return s.link.Len()
}
