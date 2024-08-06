package logger

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

var levelColor map[interfaces.LogLevel]string = map[interfaces.LogLevel]string{
	interfaces.DEBUG:    "\033[35m",
	interfaces.INFO:     "\033[36m",
	interfaces.WARNING:  "\033[33m",
	interfaces.MESSAGE:  "\033[32m",
	interfaces.ERROR:    "\033[31m",
	interfaces.CRITICAL: "\033[31m",
}

type Std struct {
	Parent interfaces.ILogger
	Config StdConfig
	queue  chan string
	done   chan interface{}
}

type StdConfig struct {
	Level interfaces.LogLevel `yaml:"level"`
}

func NewStd(cfg StdConfig, parent interfaces.ILogger) *Std {
	return &Std{
		Config: cfg,
		Parent: parent,
	}
}

func (t *Std) Init() error {
	t.queue = make(chan string, 1000)
	t.done = make(chan interface{})

	t.run()

	if t.Parent != nil {
		return t.Parent.Init()
	}

	return nil
}

func (t *Std) Stop() error {
	if t.Parent != nil {
		return t.Parent.Stop()
	}

	close(t.done)
	time.Sleep(time.Second * 1)
	close(t.queue)

	return nil
}

func (t *Std) run() {
	go func() {
		for {
			select {
			case m := <-t.queue:
				fmt.Println(m)

			case <-t.done:
				time.Sleep(time.Millisecond * 100)
				return
			}
		}
	}()
}

func (t *Std) SetLevel(level interfaces.LogLevel) {
	t.Config.Level = level

	if t.Parent != nil {
		t.Parent.SetLevel(level)
	}
}

func (t *Std) Push(level interfaces.LogLevel, msg string, toParent bool) {
	fmt.Printf("%s%s\u001B[0m [\033[37m%s\033[0m] %s\n", levelColor[level], level.String(), time.Now().Format("2006-01-02 15:04.05"), msg)

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
