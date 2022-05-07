package list

import (
	"container/list"
	"sync"
)

type List interface {
	Push(v interface{})
	Pop() interface{}
	Head() interface{}
	Tail() interface{}
	Len() int
}

// 队列
type LinkList struct {
	values *list.List
	sync.RWMutex
}

func NewLinkList() *LinkList {
	return &LinkList{
		values: list.New(),
	}
}

func (l *LinkList) PushFront(v interface{}) {
	l.Lock()
	defer l.Unlock()
	l.values.PushFront(v)
}

func (l *LinkList) PushBack(v interface{}) {
	l.Lock()
	defer l.Unlock()
	l.values.PushBack(v)
}

func (l *LinkList) PopFront() interface{} {
	l.Lock()
	defer l.Unlock()

	if v := l.values.Front(); v != nil {
		l.values.Remove(v)
		return v.Value
	}
	return nil
}

func (l *LinkList) PopBack() interface{} {
	l.Lock()
	defer l.Unlock()

	if v := l.values.Back(); v != nil {
		l.values.Remove(v)
		return v.Value
	}
	return nil
}

func (l *LinkList) Head() interface{} {
	l.RLock()
	defer l.RUnlock()
	if v := l.values.Front(); v != nil {
		return v.Value
	}
	return nil
}

func (l *LinkList) Tail() interface{} {
	l.RLock()
	defer l.RUnlock()
	if v := l.values.Back(); v != nil {
		return v.Value
	}
	return nil
}

func (l *LinkList) Len() int {
	l.RLock()
	defer l.RUnlock()
	return l.values.Len()
}
