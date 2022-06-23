package lru

import (
	"fmt"
	"testing"

	"github.com/chuan-fu/Common/baseservice/cast"
)

func TestLru(t *testing.T) {
	l := NewLruCache(10)
	fmt.Println("=====top=====")
	top := l.Top(-1)
	for _, v := range top {
		fmt.Println(v.K(), v.V())
	}

	for i := 0; i < 25; i++ {
		l.Set(cast.ToString(i), fmt.Sprintf("Vaule:%d", i))
	}
	for i := 25; i > 18; i-- {
		l.Set(cast.ToString(i), fmt.Sprintf("Vaule2:%d", i))
	}
	fmt.Println("=====top=====")
	top = l.Top(0)
	for _, v := range top {
		fmt.Println(v.K(), v.V())
	}
}
