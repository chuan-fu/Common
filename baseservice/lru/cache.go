package lru

import (
	"container/list"
	"sync"

	"github.com/chuan-fu/Common/baseservice/cast"
)

// LRU cache实现
type LruCache struct {
	size int
	list *list.List

	cacheMap map[string]*list.Element
	sync.RWMutex
}

func NewLruCache(size int) *LruCache {
	return &LruCache{
		size:     size,
		list:     list.New(),
		cacheMap: make(map[string]*list.Element, size),
	}
}

func (l *LruCache) Set(k string, v interface{}) {
	l.Lock()
	defer l.Unlock()

	if node, ok := l.cacheMap[k]; ok { // 存在的话，节点提前到第一位，修改数据
		l.list.MoveToFront(node)
		e, _ := node.Value.(*Entry)
		e.v = v
		return
	}

	if l.list.Len() == l.size {
		back := l.list.Back()
		l.list.Remove(back)
		e, _ := back.Value.(*Entry)
		delete(l.cacheMap, e.k)
	}
	l.cacheMap[k] = l.list.PushFront(&Entry{k: k, v: v})
}

// 查询
func (l *LruCache) Get(k string) (*Entry, bool) {
	l.Lock()
	defer l.Unlock()

	if node, ok := l.cacheMap[k]; ok {
		l.list.MoveToFront(node)
		e, _ := node.Value.(*Entry)
		return e, true
	}
	return nil, false
}

func (l *LruCache) TopOne() *Entry {
	l.RLock()
	defer l.RUnlock()

	node := l.list.Front()
	if node == nil {
		return nil
	}
	e, _ := node.Value.(*Entry)
	return e
}

func (l *LruCache) Top(size int) (resp []*Entry) {
	l.RLock()
	defer l.RUnlock()

	if size <= 0 {
		size = l.list.Len()
	}
	num := min(l.list.Len(), size)

	resp = make([]*Entry, 0, num)

	node := l.list.Front()
	for i := 0; i < num; i++ {
		if node == nil {
			return
		}
		e, _ := node.Value.(*Entry)
		resp = append(resp, e)
		node = node.Next()
	}
	return
}

type Entry struct {
	k string
	v interface{}
}

func (e *Entry) K() interface{} {
	return e.k
}

func (e *Entry) V() interface{} {
	return e.v
}

func (e *Entry) StringV() string {
	return cast.ToString(e.v)
}

func (e *Entry) Int64V() (int64, error) {
	return cast.ToInt64I(e.v)
}

func (e *Entry) IntV() (int, error) {
	return cast.ToIntI(e.v)
}

func (e *Entry) Float64V() (float64, error) {
	return cast.ToFloat64I(e.v)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
