package zlog

import "github.com/rs/zerolog"

type ZlogConf struct {
	SysName  string `default:"Common"`
	LogLevel string `default:"info"`
}

func (z *ZlogConf) getLevel() zerolog.Level {
	switch z.LogLevel {
	case zerolog.LevelTraceValue:
		return zerolog.TraceLevel
	case zerolog.LevelDebugValue:
		return zerolog.DebugLevel
	case zerolog.LevelInfoValue:
		return zerolog.InfoLevel
	case zerolog.LevelWarnValue:
		return zerolog.WarnLevel
	case zerolog.LevelErrorValue:
		return zerolog.ErrorLevel
	case zerolog.LevelFatalValue:
		return zerolog.FatalLevel
	case zerolog.LevelPanicValue:
		return zerolog.PanicLevel
	default:
		return zerolog.TraceLevel
	}
}
