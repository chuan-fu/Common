package log

import "fmt"

func Trace(msg ...interface{}) {
	Logger().TraceEvent().Msg(fmt.Sprint(msg...))
}

func Tracef(format string, a ...interface{}) {
	Logger().TraceEvent().Msg(fmt.Sprintf(format, a...))
}

func Debug(msg ...interface{}) {
	Logger().DebugEvent().Msg(fmt.Sprint(msg...))
}

func Debugf(format string, a ...interface{}) {
	Logger().DebugEvent().Msg(fmt.Sprintf(format, a...))
}

func Info(msg ...interface{}) {
	Logger().InfoEvent().Msg(fmt.Sprint(msg...))
}

func Infof(format string, a ...interface{}) {
	Logger().InfoEvent().Msg(fmt.Sprintf(format, a...))
}

func Warn(msg ...interface{}) {
	Logger().WarnEvent().Msg(fmt.Sprint(msg...))
}

func Warnf(format string, a ...interface{}) {
	Logger().WarnEvent().Msg(fmt.Sprintf(format, a...))
}

func Error(msg ...interface{}) {
	Logger().ErrorEvent().Msg(fmt.Sprint(msg...))
}

func Errorf(format string, a ...interface{}) {
	Logger().ErrorEvent().Msg(fmt.Sprintf(format, a...))
}

func Fatal(msg ...interface{}) {
	Logger().FatalEvent().Msg(fmt.Sprint(msg...))
}

func Fatalf(format string, a ...interface{}) {
	Logger().FatalEvent().Msg(fmt.Sprintf(format, a...))
}

func Panic(msg ...interface{}) {
	Logger().PanicEvent().Msg(fmt.Sprint(msg...))
}

func Panicf(format string, a ...interface{}) {
	Logger().PanicEvent().Msg(fmt.Sprintf(format, a...))
}
