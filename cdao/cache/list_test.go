package cache

import (
	"fmt"
	"reflect"
	"testing"
)

func A(obj interface{}) {
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			fmt.Println(ele.Interface())
		}
	}
}

func TestA(t *testing.T) {
	a := []AA{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
	}
	A(a)
}
