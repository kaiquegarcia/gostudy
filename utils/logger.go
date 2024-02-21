package utils

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota + 1
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelPanic
)

var levelLabels = map[LogLevel]string{
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelWarning: "WARN",
	LevelError:   "ERROR",
	LevelFatal:   "FATAL",
	LevelPanic:   "PANIC",
}

func (ll LogLevel) String() string {
	if label, exists := levelLabels[ll]; exists {
		return label
	}

	return "??"
}

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(err error, format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})
}

func NewLogger(
	printer Printer,
	lowestLevelAllowed LogLevel,
) Logger {
	return &logger{
		printer:            printer,
		lowestLevelAllowed: lowestLevelAllowed,
	}
}

type logger struct {
	printer            Printer
	lowestLevelAllowed LogLevel
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args)
}

func (l *logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarning, format, args)
}

func (l *logger) Error(err error, format string, args ...interface{}) {
	l.log(
		LevelError,
		"something went wrong:\n-------- message: %s\n-------- error: %s\n",
		[]interface{}{
			fmt.Sprintf(format, args...),
			err,
		},
	)
}

func (l *logger) Fatal(format string, args ...interface{}) {
	l.log(LevelFatal, format, args)
}

func (l *logger) Panic(format string, args ...interface{}) {
	l.log(LevelPanic, format, args)
}

func (l *logger) log(level LogLevel, format string, args []interface{}) {
	if level < l.lowestLevelAllowed {
		return
	}

	l.printer.Printf(
		"%s [%s] %s\n",
		time.Now().Format(time.RFC3339),
		level,
		fmt.Sprintf(format, args...),
	)
}
