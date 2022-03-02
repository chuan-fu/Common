package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chuan-fu/Common/util"

	"github.com/rs/zerolog"
)

func newConsoleWriter() io.Writer {
	return zerolog.ConsoleWriter{
		Out:         os.Stderr,
		TimeFormat:  defaultTimeFormat,
		FormatLevel: consoleFormatLevel,
	}
}

func newJsonWriter() io.Writer {
	return os.Stdout
}

func consoleFormatLevel(i interface{}) string {
	l, ok := i.(string)
	if !ok {
		if i == nil {
			return util.Colorize("?????", util.ColorBold)
		}
		return strings.ToUpper(fmt.Sprintf("%v", i))
	}

	switch l {
	case zerolog.LevelTraceValue:
		return util.Colorize("TRACE", util.ColorMagenta)
	case zerolog.LevelDebugValue:
		return util.Colorize("DEBUG", util.ColorYellow)
	case zerolog.LevelInfoValue:
		return util.Colorize("INFO ", util.ColorGreen)
	case zerolog.LevelWarnValue:
		return util.Colorize("WARN ", util.ColorRed)
	case zerolog.LevelErrorValue:
		return util.Colorize(util.Colorize("ERROR", util.ColorRed), util.ColorBold)
	case zerolog.LevelFatalValue:
		return util.Colorize(util.Colorize("FATAL", util.ColorRed), util.ColorBold)
	case zerolog.LevelPanicValue:
		return util.Colorize(util.Colorize("PANIC", util.ColorRed), util.ColorBold)
	default:
		return util.Colorize("?????", util.ColorBold)
	}
}
