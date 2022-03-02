package log

import (
	"runtime/debug"

	"github.com/rs/zerolog"
)

type StackHook struct{}

func NewStackHook() zerolog.Hook {
	return &StackHook{}
}

func (StackHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.FatalLevel || level == zerolog.PanicLevel {
		e.Str("stack", string(debug.Stack()))
	}
}
