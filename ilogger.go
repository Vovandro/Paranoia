package goServer

import "strings"

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
	Init(app *Service) error
	Stop() error
	SetLevel(level LogLevel)
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Message(args ...interface{})
	Error(err error)
	Fatal(err error)
	Panic(err error)
}

func (t LogLevel) GetLogName() string {
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
