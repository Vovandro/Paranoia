package logger

import (
	"Paranoia/interfaces"
	"fmt"
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

func (t *Std) Debug(args ...interface{}) {
	if t.Level <= interfaces.DEBUG {
		fmt.Println(args...)

		if t.Parent != nil {
			t.Parent.Debug(args)
		}
	}
}

func (t *Std) Info(args ...interface{}) {
	if t.Level <= interfaces.INFO {
		fmt.Println(args...)

		if t.Parent != nil {
			t.Parent.Info(args)
		}
	}
}

func (t *Std) Warn(args ...interface{}) {
	if t.Level <= interfaces.WARNING {
		fmt.Println(args...)

		if t.Parent != nil {
			t.Parent.Warn(args)
		}
	}
}

func (t *Std) Message(args ...interface{}) {
	if t.Level <= interfaces.MESSAGE {
		fmt.Println(args...)

		if t.Parent != nil {
			t.Parent.Message(args)
		}
	}
}

func (t *Std) Error(err error) {
	if t.Level <= interfaces.ERROR {
		fmt.Println(err)

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *Std) Fatal(err error) {
	if t.Level <= interfaces.CRITICAL {
		fmt.Println(err)

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *Std) Panic(err error) {
	if t.Level <= interfaces.CRITICAL {
		fmt.Println(err)

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
