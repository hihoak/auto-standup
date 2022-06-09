package utils

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Log - main Logger for program
var Log Logger

// Logger - struct of Logger
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

// NewLogger - creater for Logger
func NewLogger(logLevel string) Logger {
	return Logger{
		logger: zerolog.New(os.Stdout).Level(convertLogLevelToZerolog(logLevel)),
	}
}

// Info - ...
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Debug - ...
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Warn - ...
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Fatal - ...
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// Error - ,,,
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}
