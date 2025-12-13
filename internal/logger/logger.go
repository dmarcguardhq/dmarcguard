package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(logLevel string, noColor bool) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	level := zerolog.InfoLevel
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "critical":
		level = zerolog.FatalLevel
	}

	l := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    noColor,
	}).Level(level).With().Caller().Timestamp().Logger()

	return &l
}
