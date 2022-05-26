package util

import (
	"runtime"
	"runtime/debug"

	"github.com/chuan-fu/Common/zlog"
)

func DeferFunc() {
	if e := recover(); e != nil {
		DeferFuncLog(e)
	}
}

func DeferFuncLog(e interface{}) {
	var file string
	var line int
	pc, _, _, ok := runtime.Caller(3)
	if ok {
		file, line = runtime.FuncForPC(pc).FileLine(pc) // main.(*MyStruct).foo
	}
	log.Errorf("panic at(%s:%d): %v", file, line, ToString(e))
	log.Error(BytesToString(debug.Stack()))
}

func Go(f func()) {
	go func() {
		defer DeferFunc()
		f()
	}()
}
