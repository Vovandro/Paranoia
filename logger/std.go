package logger

import (
	"Paranoia/interfaces"
	"fmt"
	"time"
)

type Std struct {
	Parent interfaces.ILogger
	Level  interfaces.LogLevel
}

func (t *Std) Init(app interfaces.IService) error {
	if t.Parent != nil {
		return t.Parent.Init(app)
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
	t.Level = level

	if t.Parent != nil {
		t.Parent.SetLevel(level)
	}
}

func (t *Std) Push(level interfaces.LogLevel, msg string, toParent bool) {
	fmt.Printf("%s [%v] %s", level.String(), time.Now(), msg)

	if toParent && t.Parent != nil {
		t.Parent.Push(level, msg, true)
	}
}

func (t *Std) Debug(args ...interface{}) {
	if t.Level <= interfaces.DEBUG {
		t.Push(interfaces.DEBUG, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Debug(args)
		}
	}
}

func (t *Std) Info(args ...interface{}) {
	if t.Level <= interfaces.INFO {
		t.Push(interfaces.INFO, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Info(args)
		}
	}
}

func (t *Std) Warn(args ...interface{}) {
	if t.Level <= interfaces.WARNING {
		t.Push(interfaces.WARNING, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Warn(args)
		}
	}
}

func (t *Std) Message(args ...interface{}) {
	if t.Level <= interfaces.MESSAGE {
		t.Push(interfaces.MESSAGE, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Message(args)
		}
	}
}

func (t *Std) Error(err error) {
	if t.Level <= interfaces.ERROR {
		t.Push(interfaces.DEBUG, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *Std) Fatal(err error) {
	if t.Level <= interfaces.CRITICAL {
		t.Push(interfaces.DEBUG, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *Std) Panic(err error) {
	if t.Level <= interfaces.CRITICAL {
		t.Push(interfaces.DEBUG, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
