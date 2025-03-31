package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

const (
	DebugLevelValue = "debug"
	InfoLevelValue  = "info"
	WarnLevelValue  = "warn"
	ErrorLevelValue = "error"
	FatalLevelValue = "fatal"
)

func (l Level) zerolog() zerolog.Level {
	switch l {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	}

	return zerolog.DebugLevel
}

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return DebugLevelValue
	case InfoLevel:
		return InfoLevelValue
	case WarnLevel:
		return WarnLevelValue
	case ErrorLevel:
		return ErrorLevelValue
	case FatalLevel:
		return FatalLevelValue
	}
	return DebugLevelValue
}

type Logger struct {
	Logger zerolog.Logger
	level  Level
}

func NewConsoleLogger(level Level) *Logger {
	lvl := level.zerolog()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(lvl).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).
		Logger()
	return &Logger{
		Logger: logger,
		level:  level,
	}
}

func (l *Logger) Trace(msg string) {
	l.Logger.Trace().Msg(msg)
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}
func (l *Logger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}
func (l *Logger) Warn(msg string) {
	l.Logger.Warn().Msg(msg)
}
func (l *Logger) Error(msg string) {
	l.Logger.Error().Msg(msg)
}
