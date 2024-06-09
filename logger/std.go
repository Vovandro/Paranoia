package logger

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type Std struct {
	Parent interfaces.ILogger
	Config StdConfig
}

type StdConfig struct {
	Level interfaces.LogLevel `yaml:"level"`
}

func NewStd(cfg StdConfig) *Std {
	return &Std{
		Config: cfg,
	}
}

func (t *Std) Init() error {
	if t.Parent != nil {
		return t.Parent.Init()
	}

	return nil
}

func (t *Std) Stop() error {
	if t.Parent != nil {
		return t.Parent.Stop()
	}

	return nil
}

func (t *Std) SetLevel(level interfaces.LogLevel) {
	t.Config.Level = level

	if t.Parent != nil {
		t.Parent.SetLevel(level)
	}
}

func (t *Std) Push(level interfaces.LogLevel, msg string, toParent bool) {
	fmt.Printf("%s [%s] %s\n", level.String(), time.Now().Format("2006-01-02 15:04.05"), msg)

	if toParent && t.Parent != nil {
		t.Parent.Push(level, msg, true)
	}
}

func (t *Std) Debug(args ...interface{}) {
	if t.Config.Level <= interfaces.DEBUG {
		t.Push(interfaces.DEBUG, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Debug(args...)
		}
	}
}

func (t *Std) Info(args ...interface{}) {
	if t.Config.Level <= interfaces.INFO {
		t.Push(interfaces.INFO, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Info(args...)
		}
	}
}

func (t *Std) Warn(args ...interface{}) {
	if t.Config.Level <= interfaces.WARNING {
		t.Push(interfaces.WARNING, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Warn(args...)
		}
	}
}

func (t *Std) Message(args ...interface{}) {
	if t.Config.Level <= interfaces.MESSAGE {
		t.Push(interfaces.MESSAGE, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Message(args...)
		}
	}
}

func (t *Std) Error(err error) {
	if t.Config.Level <= interfaces.ERROR {
		t.Push(interfaces.DEBUG, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *Std) Fatal(err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		t.Push(interfaces.DEBUG, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *Std) Panic(err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		t.Push(interfaces.DEBUG, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
