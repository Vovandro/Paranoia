package mock_log

import (
	"context"
	"fmt"
	"time"
)

var levelColor = map[LogLevel]string{
	DEBUG:    "\033[35m",
	INFO:     "\033[36m",
	WARNING:  "\033[33m",
	MESSAGE:  "\033[32m",
	ERROR:    "\033[31m",
	CRITICAL: "\033[31m",
}

type Mock struct {
	enable bool
}

func New(enable bool) *Mock {
	return &Mock{enable: enable}
}

func (t *Mock) Init(cfg map[string]interface{}) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return ""
}

func (t *Mock) Type() string {
	return "logger"
}

func (t *Mock) SetLevel(level int) {

}

func (t *Mock) push(level LogLevel, msg string) {
	if !t.enable {
		return
	}

	fmt.Printf("%s%s\u001B[0m [\033[37m%s\033[0m] %s\n", levelColor[level], level.String(), time.Now().Format("2006-01-02 15:04.05"), msg)
}

func (t *Mock) Debug(ctx context.Context, args ...interface{}) {
	t.push(DEBUG, fmt.Sprint(args...))
}

func (t *Mock) Info(ctx context.Context, args ...interface{}) {
	t.push(INFO, fmt.Sprint(args...))
}

func (t *Mock) Warn(ctx context.Context, args ...interface{}) {
	t.push(WARNING, fmt.Sprint(args...))
}

func (t *Mock) Message(ctx context.Context, args ...interface{}) {
	t.push(MESSAGE, fmt.Sprint(args...))
}

func (t *Mock) Error(ctx context.Context, err error) {
	t.push(ERROR, err.Error())
}

func (t *Mock) Fatal(ctx context.Context, err error) {
	t.push(CRITICAL, err.Error())
}

func (t *Mock) Panic(ctx context.Context, err error) {
	t.push(CRITICAL, err.Error())
}

func (t *Mock) Parent() interface{} {
	return nil
}

func (t *Mock) SetParent(parent interface{}) {

}
