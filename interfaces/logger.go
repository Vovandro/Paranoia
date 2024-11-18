package interfaces

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
	Init(cfg IConfig) error
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

func (t LogLevel) String() string {
	switch t {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case MESSAGE:
		return "MESSAGE"
	case WARNING:
		return "WARNING"
	case CRITICAL:
		return "CRITICAL"

	default:
		return ""
	}
}

func GetLogLevel(str string) LogLevel {
	switch strings.ToUpper(str) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "MESSAGE":
		return MESSAGE
	case "WARNING":
		return WARNING
	case "CRITICAL":
		return CRITICAL

	default:
		return DEBUG
	}
}
