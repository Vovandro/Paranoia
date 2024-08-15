package logger

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"os"
	"time"
)

type File struct {
	Parent interfaces.ILogger
	Config FileConfig
	queue  chan string
	done   chan interface{}
	f      *os.File
}

type FileConfig struct {
	Level interfaces.LogLevel `yaml:"level"`
	FName string              `yaml:"filename"`
}

func NewFile(cfg FileConfig, parent interfaces.ILogger) *File {
	return &File{
		Config: cfg,
		Parent: parent,
	}
}

func (t *File) Init(cfg interfaces.IConfig) error {
	t.queue = make(chan string, 1000)
	t.done = make(chan interface{})

	_, err := os.Stat("./log")

	if os.IsNotExist(err) {
		_ = os.Mkdir("./log", 0766)
	}

	t.f, err = os.OpenFile(
		fmt.Sprintf("./log/%s_%s.log", t.Config.FName, time.Now().Format("2006_01_02")),
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666)

	if err != nil {
		fmt.Println("error open log file")
		return err
	}

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

func (t *File) Stop() error {
	if t.Parent != nil {
		return t.Parent.Stop()
	}

	close(t.done)
	time.Sleep(time.Second * 1)
	close(t.queue)

	return nil
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
					fmt.Sprintf("./log/%s_%s.log", t.Config.FName, time.Now().Format("2006_01_02")),
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

func (t *File) SetLevel(level interfaces.LogLevel) {
	t.Config.Level = level

	if t.Parent != nil {
		t.Parent.SetLevel(level)
	}
}

func (t *File) push(level interfaces.LogLevel, msg string) {
	t.queue <- fmt.Sprintf("%s [%s] %s\n", level.String(), time.Now().Format("2006-01-02 15:04.05"), msg)
}

func (t *File) Debug(args ...interface{}) {
	if t.Config.Level <= interfaces.DEBUG {
		t.push(interfaces.DEBUG, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Debug(args...)
		}
	}
}

func (t *File) Info(args ...interface{}) {
	if t.Config.Level <= interfaces.INFO {
		t.push(interfaces.INFO, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Info(args...)
		}
	}
}

func (t *File) Warn(args ...interface{}) {
	if t.Config.Level <= interfaces.WARNING {
		t.push(interfaces.WARNING, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Warn(args...)
		}
	}
}

func (t *File) Message(args ...interface{}) {
	if t.Config.Level <= interfaces.MESSAGE {
		t.push(interfaces.MESSAGE, fmt.Sprint(args...))

		if t.Parent != nil {
			t.Parent.Message(args...)
		}
	}
}

func (t *File) Error(err error) {
	if t.Config.Level <= interfaces.ERROR {
		t.push(interfaces.ERROR, err.Error())

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *File) Fatal(err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		t.push(interfaces.CRITICAL, err.Error())

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *File) Panic(err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		t.push(interfaces.CRITICAL, err.Error())

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
