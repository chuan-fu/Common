package log

import (
	"os"

	"github.com/rs/zerolog"
)

const (
	skipFrameCount = 3
	sysNameKey     = "sysname"
)

var zlogger *logger

func init() {
	ReloadLogger(ZlogConf{
		SysName:  "Common",
		LogLevel: zerolog.LevelTraceValue,
	})
}

func ReloadLogger(conf ZlogConf) {
	zlogger = NewLogger(conf)
}

func NewLogger(conf ZlogConf) *logger {
	zerolog.SetGlobalLevel(conf.getLevel()) // 修改日志等级
	zerolog.CallerFieldName = "gofile"
	zerolog.MessageFieldName = "msg"
	l := &logger{
		log: zerolog.New(os.Stderr).
			With().Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger().
			Hook(NewStackHook()), // 添加stack钩子
		sysName: conf.SysName,
	}
	return l
}

type logger struct {
	log     zerolog.Logger
	sysName string
}

func Logger() *logger {
	return zlogger
}

func (l *logger) TraceEvent() *zerolog.Event {
	return l.log.Trace().Str(sysNameKey, l.sysName)
}

func (l *logger) DebugEvent() *zerolog.Event {
	return l.log.Debug().Str(sysNameKey, l.sysName)
}

func (l *logger) InfoEvent() *zerolog.Event {
	return l.log.Info().Str(sysNameKey, l.sysName)
}

func (l *logger) WarnEvent() *zerolog.Event {
	return l.log.Warn().Str(sysNameKey, l.sysName)
}

func (l *logger) ErrorEvent() *zerolog.Event {
	return l.log.Error().Str(sysNameKey, l.sysName)
}

func (l *logger) FatalEvent() *zerolog.Event {
	return l.log.Fatal().Str(sysNameKey, l.sysName)
}

func (l *logger) PanicEvent() *zerolog.Event {
	return l.log.Panic().Str(sysNameKey, l.sysName)
}
