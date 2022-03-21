package log

import (
	"runtime/debug"
	"unsafe"

	"github.com/rs/zerolog"
)

type StackHook struct{}

func NewStackHook() zerolog.Hook {
	return &StackHook{}
}

func (StackHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.FatalLevel || level == zerolog.PanicLevel {
		e.Str(stackKey, *BytesToString(debug.Stack()))
	}
}

// 该函数在util包中有
// 但是log不能依赖于任何内部包
// 不然可能出现循环引用
func BytesToString(b []byte) *string {
	return (*string)(unsafe.Pointer(&b))
}
