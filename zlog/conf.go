package log

import "github.com/rs/zerolog"

/*
type Config struct {
	Log log.ZlogConf `json:"log" yaml:"log"`
}

log:
	sysName: Common
	logLevel: info
	encoding: json
	noColor: true

*/

type ZlogConf struct {
	SysName  string `default:"Common" json:"sysName" yaml:"sysName"`
	LogLevel string `default:"trace" json:"logLevel" yaml:"logLevel"`
	Encoding string `default:"json" json:"encoding" yaml:"encoding"`
	NoColor  bool   `default:"true" json:"noColor" yaml:"noColor"`
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

const (
	EncodingJson    = "json"
	EncodingConsole = "console"
)

func (z *ZlogConf) getEncoding() string {
	switch z.Encoding {
	case EncodingJson, EncodingConsole:
		return z.Encoding
	default:
		return EncodingJson
	}
}
