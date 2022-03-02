package log

import (
	"fmt"
	"io"
	"os"
	"strings"

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
				return Colorize("?????", ColorBold, noColor)
			}
			return strings.ToUpper(fmt.Sprintf("%v", i))
		}

		switch l {
		case zerolog.LevelTraceValue:
			return Colorize("TRACE", ColorMagenta, noColor)
		case zerolog.LevelDebugValue:
			return Colorize("DEBUG", ColorYellow, noColor)
		case zerolog.LevelInfoValue:
			return Colorize("INFO ", ColorGreen, noColor)
		case zerolog.LevelWarnValue:
			return Colorize("WARN ", ColorRed, noColor)
		case zerolog.LevelErrorValue:
			return Colorize(Colorize("ERROR", ColorRed, noColor), ColorBold, noColor)
		case zerolog.LevelFatalValue:
			return Colorize(Colorize("FATAL", ColorRed, noColor), ColorBold, noColor)
		case zerolog.LevelPanicValue:
			return Colorize(Colorize("PANIC", ColorRed, noColor), ColorBold, noColor)
		default:
			return Colorize("?????", ColorBold, noColor)
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
			s = Colorize(s, ColorBold, noColor) + Colorize(" >", ColorCyan, noColor)
		}
		return s
	}
}
