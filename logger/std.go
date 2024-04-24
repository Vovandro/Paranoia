package logger

import (
	"Paranoia"
	"fmt"
)

type Std struct {
	Parent Paranoia.ILogger
	Level  Paranoia.LogLevel
}

func (t *Std) Init(app *Paranoia.Service) error {
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

func (t *Std) SetLevel(level Paranoia.LogLevel) {
	t.Level = level

	if t.Parent != nil {
		t.Parent.SetLevel(level)
	}
}

func (t *Std) Debug(args ...interface{}) {
	if t.Level <= Paranoia.DEBUG {
		fmt.Println(args)

		if t.Parent != nil {
			t.Parent.Debug(args)
		}
	}
}

func (t *Std) Info(args ...interface{}) {
	if t.Level <= Paranoia.INFO {
		fmt.Println(args)

		if t.Parent != nil {
			t.Parent.Info(args)
		}
	}
}

func (t *Std) Warn(args ...interface{}) {
	if t.Level <= Paranoia.WARNING {
		fmt.Println(args)

		if t.Parent != nil {
			t.Parent.Warn(args)
		}
	}
}

func (t *Std) Message(args ...interface{}) {
	if t.Level <= Paranoia.MESSAGE {
		fmt.Println(args)

		if t.Parent != nil {
			t.Parent.Message(args)
		}
	}
}

func (t *Std) Error(err error) {
	if t.Level <= Paranoia.ERROR {
		fmt.Println(err)

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *Std) Fatal(err error) {
	if t.Level <= Paranoia.CRITICAL {
		fmt.Println(err)

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *Std) Panic(err error) {
	if t.Level <= Paranoia.CRITICAL {
		fmt.Println(err)

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
