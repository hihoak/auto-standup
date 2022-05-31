package utils

import (
	"github.com/rs/zerolog"
	"os"
	"strings"
)

var Log Logger

type Logger struct {
	logger zerolog.Logger
}

func convertLogLevelToZerolog(logLevel string) zerolog.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}

func NewLogger(logLevel string) Logger {
	return Logger{
		logger: zerolog.New(os.Stdout).Level(convertLogLevelToZerolog(logLevel)),
	}
}

func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}
