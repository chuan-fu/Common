package util

import (
	"fmt"
	"testing"
)

func TestPanic(t *testing.T) {
	panic(nil)
}

func TestPanicWithDefer(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	panic(nil)
}

func TestDeferFunc(t *testing.T) {
	a := "aaa"
	defer DeferFunc()
	panic(a)
}

func TestDeferFuncLog(t *testing.T) {
	a := "aaa"
	defer func() {
		if e := recover(); e != nil {
			DeferHandle(e)
		}
	}()
	panic(a)
}
