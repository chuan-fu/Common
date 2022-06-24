package util

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/chuan-fu/Common/baseservice/stringx"

	"github.com/chuan-fu/Common/zlog"
)

const (
	SkipFour = 4
	SkipFive = 5
	SkipSix  = 6
)

var DeferFunc = DeferFuncWithSkip(SkipFour)

func DeferFuncWithSkip(skip int) func() {
	return func() {
		if e := recover(); e != nil { // panic(nil)会被捕获到，但是过不了 e!=nil 的判断
			DeferHandleWithSkip(e, skip)
		}
	}
}

// panic处理函数
func DeferHandle(e interface{}) {
	DeferHandleWithSkip(e, SkipFive)
}

// panic处理函数
func DeferHandleWithSkip(e interface{}, skip int) {
	if err := NewPanicErrorWithSkip(e, skip); err != nil {
		log.Error(err.Error())
	}
}

func Go(f func()) {
	go func() {
		defer DeferFunc()
		f()
	}()
}

type PanicError struct {
	Value, Stack string
}

func (p *PanicError) Error() string {
	return fmt.Sprintf("%s\nstack: %s", p.Value, p.Stack)
}

func NewPanicError(e interface{}) error {
	return NewPanicErrorWithSkip(e, SkipFour)
}

func NewPanicErrorWithSkip(e interface{}, skip int) error {
	var file string
	var line int
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		file, line = runtime.FuncForPC(pc).FileLine(pc) // main.(*MyStruct).foo
	}
	return &PanicError{
		Value: fmt.Sprintf("panic at(%s:%d): %v", file, line, e),
		Stack: stringx.BytesToString(debug.Stack()),
	}
}
