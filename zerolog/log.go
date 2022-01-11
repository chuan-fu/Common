package log

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.CallerSkipFrameCount = 3
	log.Logger = log.With().Caller().Logger()
}

func Info(msg ...string) {
	log.Info().Msg(strings.Join(msg, " "))
}

func Error(err error) {
	log.Error().Msg(err.Error())
}
