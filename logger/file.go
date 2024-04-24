package logger

import (
	"Paranoia/interfaces"
	"fmt"
	"os"
	"time"
)

type File struct {
	Parent interfaces.ILogger
	Level  interfaces.LogLevel
	FName  string
	queue  chan string
	done   chan interface{}
	f      *os.File
}

func (t *File) Init(app interfaces.IService) error {
	t.queue = make(chan string, 1000)
	t.done = make(chan interface{})

	_, err := os.Stat("./log")

	if os.IsNotExist(err) {
		_ = os.Mkdir("./log", 0766)
	}

	t.f, err = os.OpenFile(
		fmt.Sprintf("./log/%s%s.log", t.FName, time.Now().Format("2006_01_02")),
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666)

	if err != nil {
		fmt.Println("error open log file")
		return err
	}

	t.run(t.done)

	if t.Parent != nil {
		return t.Parent.Init(app)
	}

	return nil
}

func (t *File) Stop() error {
	close(t.done)
	close(t.queue)

	return nil
}

func (t *File) run(done chan interface{}) {
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
					fmt.Sprintf("./log/%s%s.log", t.FName, time.Now().Format("2006_01_02")),
					os.O_WRONLY|os.O_APPEND|os.O_CREATE,
					0666)

				if err != nil {
					fmt.Println("error open log file")
				}

				timeNow = time.Now()
				seconds = 24*60*60 - (timeNow.Hour()*60*60 + timeNow.Minute()*60 + timeNow.Second())
				timerRecreate.Reset(time.Second * time.Duration(seconds))

			case <-done:
				time.Sleep(time.Second * 1)
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
	t.Level = level

	if t.Parent != nil {
		t.Parent.SetLevel(level)
	}
}

func (t *File) Push(level interfaces.LogLevel, msg string, toParent bool) {
	t.queue <- fmt.Sprintf("%s [%v] %s", level.String(), time.Now(), msg)

	if toParent && t.Parent != nil {
		t.Parent.Push(level, msg, true)
	}
}

func (t *File) Debug(args ...interface{}) {
	if t.Level <= interfaces.DEBUG {
		t.Push(interfaces.DEBUG, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Debug(args)
		}
	}
}

func (t *File) Info(args ...interface{}) {
	if t.Level <= interfaces.INFO {
		t.Push(interfaces.INFO, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Info(args)
		}
	}
}

func (t *File) Warn(args ...interface{}) {
	if t.Level <= interfaces.WARNING {
		t.Push(interfaces.WARNING, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Warn(args)
		}
	}
}

func (t *File) Message(args ...interface{}) {
	if t.Level <= interfaces.MESSAGE {
		t.Push(interfaces.MESSAGE, fmt.Sprint(args...), false)

		if t.Parent != nil {
			t.Parent.Message(args)
		}
	}
}

func (t *File) Error(err error) {
	if t.Level <= interfaces.ERROR {
		t.Push(interfaces.ERROR, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *File) Fatal(err error) {
	if t.Level <= interfaces.CRITICAL {
		t.Push(interfaces.CRITICAL, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *File) Panic(err error) {
	if t.Level <= interfaces.CRITICAL {
		t.Push(interfaces.CRITICAL, err.Error(), false)

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
