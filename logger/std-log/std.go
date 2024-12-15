package std_log

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

var levelColor = map[interfaces.LogLevel]string{
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

func (t *Std) Init(cfg interfaces.IConfig) error {
	t.queue = make(chan string, 1000)
	t.done = make(chan interface{})

	if cfg != nil {
		l := cfg.GetString("LOG_LEVEL", "")

		if l != "" {
			t.Config.Level = interfaces.GetLogLevel(l)
		}
	}

	t.run()

	if t.Parent != nil {
		return t.Parent.Init(cfg)
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

func (t *Std) push(level interfaces.LogLevel, msg string) {
	t.queue <- fmt.Sprintf("%s%s\u001B[0m [\033[37m%s\033[0m] %s\n", levelColor[level], level.String(), time.Now().Format("2006-01-02 15:04.05"), msg)
}

func (t *Std) Debug(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.DEBUG {
		t.push(interfaces.DEBUG, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Debug(ctx, args...)
		}
	}
}

func (t *Std) Info(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.INFO {
		t.push(interfaces.INFO, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Info(ctx, args...)
		}
	}
}

func (t *Std) Warn(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.WARNING {
		t.push(interfaces.WARNING, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Warn(ctx, args...)
		}
	}
}

func (t *Std) Message(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.MESSAGE {
		t.push(interfaces.MESSAGE, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Message(ctx, args...)
		}
	}
}

func (t *Std) Error(ctx context.Context, err error) {
	if t.Config.Level <= interfaces.ERROR {
		t.push(interfaces.ERROR, err.Error())

		if t.Parent != nil {
			t.Parent.Error(ctx, err)
		}
	}
}

func (t *Std) Fatal(ctx context.Context, err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		t.push(interfaces.CRITICAL, err.Error())

		if t.Parent != nil {
			t.Parent.Fatal(ctx, err)
		}
	}
}

func (t *Std) Panic(ctx context.Context, err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		t.push(interfaces.CRITICAL, err.Error())

		if t.Parent != nil {
			t.Parent.Panic(ctx, err)
		}
	}
}
