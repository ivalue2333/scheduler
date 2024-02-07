package scheduler

import (
	"context"
	"fmt"
)

// LogLevel log level
type LogLevel int

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
	// Debug info log level
	Debug
)

type Logger interface {
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
}

func newDefaultLogger() *defaultLogger {
	return &defaultLogger{
		LogLevel: Info,
	}
}

type defaultLogger struct {
	LogLevel LogLevel
}

func (l *defaultLogger) Debug(ctx context.Context, format string, a ...interface{}) {
	if l.LogLevel >= Debug {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

func (l *defaultLogger) Info(ctx context.Context, format string, a ...interface{}) {
	if l.LogLevel >= Info {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

func (l *defaultLogger) Warn(ctx context.Context, format string, a ...interface{}) {
	if l.LogLevel >= Warn {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

func (l *defaultLogger) Error(ctx context.Context, format string, a ...interface{}) {
	if l.LogLevel >= Error {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}
