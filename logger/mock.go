package logger

import (
	"Paranoia"
	"fmt"
)

type Mock struct {
}

func (t *Mock) Init(app *Paranoia.Service) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) SetLevel(level Paranoia.LogLevel) {}
func (t *Mock) Debug(args ...interface{})        {}
func (t *Mock) Info(args ...interface{})         {}
func (t *Mock) Warn(args ...interface{})         {}
func (t *Mock) Message(args ...interface{})      {}
func (t *Mock) Error(err error) {
	fmt.Println(err)
}
func (t *Mock) Fatal(err error) {
	fmt.Println(err)
}
func (t *Mock) Panic(err error) {
	fmt.Println(err)
}
