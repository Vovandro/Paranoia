package file_log

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/devpro_studio/go_utils/decode"
	"os"
	"time"
)

type File struct {
	name   string
	parent ILogger
	config Config
	queue  chan string
	done   chan interface{}
	f      *os.File
}

type Config struct {
	Level  LogLevel `yaml:"level"`
	FName  string   `yaml:"filename"`
	Enable bool     `yaml:"enable"`
}

func New(name string) *File {
	return &File{
		name: name,
	}
}

func (t *File) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.FName == "" {
		return errors.New("filename is required")
	}

	t.queue = make(chan string, 1000)
	t.done = make(chan interface{})

	_, err = os.Stat("./log")

	if os.IsNotExist(err) {
		_ = os.Mkdir("./log", 0766)
	}

	t.f, err = os.OpenFile(
		fmt.Sprintf("./log/%s_%s.log", t.config.FName, time.Now().Format("2006_01_02")),
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666)

	if err != nil {
		fmt.Println("error open log file")
		return err
	}

	t.run()

	return nil
}

func (t *File) Stop() error {
	if t.parent != nil {
		return t.parent.Stop()
	}

	close(t.done)
	time.Sleep(time.Second * 1)
	close(t.queue)

	return nil
}

func (t *File) Name() string {
	return t.name
}

func (t *File) Type() string {
	return "logger"
}

func (t *File) run() {
	go func() {
		defer t.f.Close()

		timeNow := time.Now()
		seconds := 24*60*60 - (timeNow.Hour()*60*60 + timeNow.Minute()*60 + timeNow.Second())
		timerRecreate := time.NewTimer(time.Second * time.Duration(seconds))

		for {
			select {
			case m := <-t.queue:
				t.write(m)

			case <-timerRecreate.C:
				err := t.f.Close()

				if err != nil {
					fmt.Println("error close log file")
				}

				t.f, err = os.OpenFile(
					fmt.Sprintf("./log/%s_%s.log", t.config.FName, time.Now().Format("2006_01_02")),
					os.O_WRONLY|os.O_APPEND|os.O_CREATE,
					0666)

				if err != nil {
					fmt.Println("error open log file")
				}

				timeNow = time.Now()
				seconds = 24*60*60 - (timeNow.Hour()*60*60 + timeNow.Minute()*60 + timeNow.Second())
				timerRecreate.Reset(time.Second * time.Duration(seconds))

			case <-t.done:
				time.Sleep(time.Millisecond * 100)
				return
			}
		}
	}()
}

func (t *File) write(m string) {
	_, err := t.f.Write([]byte(m))
	if err != nil {
		fmt.Println("error write log file")
	}
}

func (t *File) SetLevel(level LogLevel) {
	t.config.Level = level

	if t.parent != nil {
		t.parent.SetLevel(level)
	}
}

func (t *File) push(level LogLevel, msg string) {
	if t.config.Enable {
		t.queue <- fmt.Sprintf("%s [%s] %s\n", level.String(), time.Now().Format("2006-01-02 15:04.05"), msg)
	}
}

func (t *File) Debug(ctx context.Context, args ...interface{}) {
	if t.config.Level <= DEBUG {
		t.push(DEBUG, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Debug(ctx, args...)
		}
	}
}

func (t *File) Info(ctx context.Context, args ...interface{}) {
	if t.config.Level <= INFO {
		t.push(INFO, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Info(ctx, args...)
		}
	}
}

func (t *File) Warn(ctx context.Context, args ...interface{}) {
	if t.config.Level <= WARNING {
		t.push(WARNING, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Warn(ctx, args...)
		}
	}
}

func (t *File) Message(ctx context.Context, args ...interface{}) {
	if t.config.Level <= MESSAGE {
		t.push(MESSAGE, fmt.Sprint(args...))

		if t.parent != nil {
			t.parent.Message(ctx, args...)
		}
	}
}

func (t *File) Error(ctx context.Context, err error) {
	if t.config.Level <= ERROR {
		t.push(ERROR, err.Error())

		if t.parent != nil {
			t.parent.Error(ctx, err)
		}
	}
}

func (t *File) Fatal(ctx context.Context, err error) {
	if t.config.Level <= CRITICAL {
		t.push(CRITICAL, err.Error())

		if t.parent != nil {
			t.parent.Fatal(ctx, err)
		}
	}
}

func (t *File) Panic(ctx context.Context, err error) {
	if t.config.Level <= CRITICAL {
		t.push(CRITICAL, err.Error())

		if t.parent != nil {
			t.parent.Panic(ctx, err)
		}
	}
}

func (t *File) Parent() interface{} {
	return t.parent
}

func (t *File) SetParent(parent interface{}) {
	t.parent = parent.(ILogger)
}
