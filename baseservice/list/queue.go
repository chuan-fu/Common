package list

// 队列
type Queue struct {
	link *LinkList
}

func NewQueue() List {
	return &Queue{NewLinkList()}
}

func (s *Queue) Push(v interface{}) {
	s.link.PushBack(v)
}

func (s *Queue) Pop() interface{} {
	return s.link.PopFront()
}

func (s *Queue) Head() interface{} {
	return s.link.Head()
}

func (s *Queue) Tail() interface{} {
	return s.link.Tail()
}

func (s *Queue) Len() int {
	return s.link.Len()
}
