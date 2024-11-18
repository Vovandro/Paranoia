package logger

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
)

type Mock struct {
}

func (t *Mock) Init(cfg interfaces.IConfig) error {
	return nil
}

func NewMock() *Mock {
	return &Mock{}
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) SetLevel(level interfaces.LogLevel)               {}
func (t *Mock) Debug(ctx context.Context, args ...interface{})   {}
func (t *Mock) Info(ctx context.Context, args ...interface{})    {}
func (t *Mock) Warn(ctx context.Context, args ...interface{})    {}
func (t *Mock) Message(ctx context.Context, args ...interface{}) {}
func (t *Mock) Error(ctx context.Context, err error) {
	fmt.Println(err)
}
func (t *Mock) Fatal(ctx context.Context, err error) {
	fmt.Println(err)
}
func (t *Mock) Panic(ctx context.Context, err error) {
	fmt.Println(err)
}
