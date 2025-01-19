package std_log

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/go_utils/decode"
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

type Std struct {
	name   string
	parent ILogger
	config Config
	queue  chan string
	done   chan interface{}
}

type Config struct {
	Level  LogLevel `yaml:"level"`
	Enable bool     `yaml:"enable"`
}

func New(name string) *Std {
	return &Std{
		name: name,
	}
}

func (t *Std) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	t.queue = make(chan string, 1000)
	t.done = make(chan interface{})

	t.run()

	return nil
}

func (t *Std) Stop() error {
	if t.parent != nil {
		return t.parent.Stop()
	}

	close(t.done)
	time.Sleep(time.Second * 1)
	close(t.queue)

	return nil
}

func (t *Std) Name() string {
	return t.name
}

func (t *Std) Type() string {
	return "logger"
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

func (t *Std) SetLevel(level LogLevel) {
	t.config.Level = level

	if t.parent != nil {
		t.parent.SetLevel(level)
	}
}

func (t *Std) push(level LogLevel, msg string) {
	if t.config.Enable {
		t.queue <- fmt.Sprintf("%s%s\u001B[0m [\033[37m%s\033[0m] %s\n", levelColor[level], level.String(), time.Now().Format("2006-01-02 15:04.05"), msg)
	}
}

func (t *Std) Debug(ctx context.Context, args ...interface{}) {
	if t.config.Level <= DEBUG {
		t.push(DEBUG, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Debug(ctx, args...)
		}
	}
}

func (t *Std) Info(ctx context.Context, args ...interface{}) {
	if t.config.Level <= INFO {
		t.push(INFO, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Info(ctx, args...)
		}
	}
}

func (t *Std) Warn(ctx context.Context, args ...interface{}) {
	if t.config.Level <= WARNING {
		t.push(WARNING, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Warn(ctx, args...)
		}
	}
}

func (t *Std) Message(ctx context.Context, args ...interface{}) {
	if t.config.Level <= MESSAGE {
		t.push(MESSAGE, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Message(ctx, args...)
		}
	}
}

func (t *Std) Error(ctx context.Context, err error) {
	if t.config.Level <= ERROR {
		t.push(ERROR, err.Error())

		if t.parent != nil {
			t.parent.Error(ctx, err)
		}
	}
}

func (t *Std) Fatal(ctx context.Context, err error) {
	if t.config.Level <= CRITICAL {
		t.push(CRITICAL, err.Error())

		if t.parent != nil {
			t.parent.Fatal(ctx, err)
		}
	}
}

func (t *Std) Panic(ctx context.Context, err error) {
	if t.config.Level <= CRITICAL {
		t.push(CRITICAL, err.Error())

		if t.parent != nil {
			t.parent.Panic(ctx, err)
		}
	}
}

func (t *Std) Parent() interface{} {
	return t.parent
}

func (t *Std) SetParent(parent interface{}) {
	t.parent = parent.(ILogger)
}
