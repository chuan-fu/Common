package util

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

func DeferFunc() {
	if e := recover(); e != nil {
		var funcName string
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name() // main.(*MyStruct).foo
		}
		err := errors.New(fmt.Sprintf("panic at(%v): %v", funcName, e))
		log.Error(err)
		log.Errorf(string(debug.Stack()))
	}
}

func Go(f func()) {
	go func() {
		defer DeferFunc()
		f()
	}()
}
