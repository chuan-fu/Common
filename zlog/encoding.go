package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chuan-fu/Common/util"

	"github.com/rs/zerolog"
)

func newConsoleWriter(noColor bool) io.Writer {
	return zerolog.ConsoleWriter{
		Out:          os.Stderr,
		NoColor:      noColor,
		TimeFormat:   defaultTimeFormat,
		FormatLevel:  consoleFormatLevel(noColor),
		FormatCaller: consoleFormatCaller(noColor),
	}
}

func newJsonWriter() io.Writer {
	return os.Stdout
}

func consoleFormatLevel(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		l, ok := i.(string)
		if !ok {
			if i == nil {
				return util.Colorize("?????", util.ColorBold, noColor)
			}
			return strings.ToUpper(fmt.Sprintf("%v", i))
		}

		switch l {
		case zerolog.LevelTraceValue:
			return util.Colorize("TRACE", util.ColorMagenta, noColor)
		case zerolog.LevelDebugValue:
			return util.Colorize("DEBUG", util.ColorYellow, noColor)
		case zerolog.LevelInfoValue:
			return util.Colorize("INFO ", util.ColorGreen, noColor)
		case zerolog.LevelWarnValue:
			return util.Colorize("WARN ", util.ColorRed, noColor)
		case zerolog.LevelErrorValue:
			return util.Colorize(util.Colorize("ERROR", util.ColorRed, noColor), util.ColorBold, noColor)
		case zerolog.LevelFatalValue:
			return util.Colorize(util.Colorize("FATAL", util.ColorRed, noColor), util.ColorBold, noColor)
		case zerolog.LevelPanicValue:
			return util.Colorize(util.Colorize("PANIC", util.ColorRed, noColor), util.ColorBold, noColor)
		default:
			return util.Colorize("?????", util.ColorBold, noColor)
		}
	}
}

func consoleFormatCaller(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		if i == nil {
			return ""
		}
		s, ok := i.(string)
		if !ok {
			s = fmt.Sprintf("%v", i)
		}

		if s != "" {
			s = util.Colorize(s, util.ColorBold, noColor) + util.Colorize(" >", util.ColorCyan, noColor)
		}
		return s
	}
}
