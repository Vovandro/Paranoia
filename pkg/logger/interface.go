package logger

import (
	"context"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	MESSAGE
	ERROR
	CRITICAL
)

type ILogger interface {
	Init(map[string]interface{}) error
	Stop() error
	SetLevel(level LogLevel)
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Message(ctx context.Context, args ...interface{})
	Error(ctx context.Context, err error)
	Fatal(ctx context.Context, err error)
	Panic(ctx context.Context, err error)
}

func (t *LogLevel) String() string {
	switch *t {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case MESSAGE:
		return "MESSAGE"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"

	default:
		return ""
	}
}

func (t *LogLevel) Parse(str string) {
	switch strings.ToUpper(str) {
	case "DEBUG":
		*t = DEBUG
	case "INFO":
		*t = INFO
	case "MESSAGE":
		*t = MESSAGE
	case "WARNING":
		*t = WARNING
	case "ERROR":
		*t = ERROR
	case "CRITICAL":
		*t = CRITICAL

	default:
		*t = DEBUG
	}
}
