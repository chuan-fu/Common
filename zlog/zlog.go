package log

import (
	"io"

	"github.com/rs/zerolog"
)

const (
	skipFrameCount    = 3
	sysNameKey        = "sysname"
	stackKey          = "stack"
	defaultTimeFormat = "2006-01-02 15:04:05"
	defaultSysName    = "Common"

	defaultCallerName  = "gofile"
	defaultMessageName = "msg"
)

var zlogger *logger

func init() {
	ReloadLogger(ZlogConf{
		SysName:  defaultSysName,
		LogLevel: zerolog.LevelTraceValue,
		Encoding: EncodingJson,
		NoColor:  true,
	})
}

func ReloadLogger(conf ZlogConf, ws ...io.Writer) {
	zlogger = NewLogger(conf, ws...)
}

func NewLogger(conf ZlogConf, ws ...io.Writer) *logger {
	zerolog.SetGlobalLevel(conf.getLevel()) // 修改日志等级
	zerolog.CallerFieldName = defaultCallerName
	zerolog.MessageFieldName = defaultMessageName
	zerolog.TimeFieldFormat = defaultTimeFormat
	l := &logger{
		log: newEncodingLogger(&conf, ws...).With().Str(sysNameKey, conf.SysName).Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger().
			Hook(NewStackHook()), // 添加stack钩子
		sysName: conf.SysName,
	}
	return l
}

func newEncodingLogger(c *ZlogConf, ws ...io.Writer) zerolog.Logger {
	var w io.Writer
	switch c.getEncoding() {
	case EncodingConsole:
		w = newConsoleWriter(c.NoColor)
	default:
		w = newJsonWriter()
	}
	if len(ws) > 0 {
		return zerolog.New(zerolog.MultiLevelWriter(append(ws, w)...))
	}
	return zerolog.New(w)
}

func SetTimeFieldFormat(data string) {
	zerolog.TimeFieldFormat = data
}

type logger struct {
	log     zerolog.Logger
	sysName string
}

func Logger() *logger {
	return zlogger
}

func (l *logger) TraceEvent() *zerolog.Event {
	return l.log.Trace()
}

func (l *logger) DebugEvent() *zerolog.Event {
	return l.log.Debug()
}

func (l *logger) InfoEvent() *zerolog.Event {
	return l.log.Info()
}

func (l *logger) WarnEvent() *zerolog.Event {
	return l.log.Warn()
}

func (l *logger) ErrorEvent() *zerolog.Event {
	return l.log.Error()
}

func (l *logger) FatalEvent() *zerolog.Event {
	return l.log.Fatal()
}

func (l *logger) PanicEvent() *zerolog.Event {
	return l.log.Panic()
}
