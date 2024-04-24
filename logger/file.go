package logger

import (
	"Paranoia/interfaces"
	"fmt"
	"os"
	"time"
)

type File struct {
	FName  string
	Parent interfaces.ILogger
	Level  interfaces.LogLevel
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
		fmt.Sprintf("./log/%s.log", time.Now().Format("2006_01_02")),
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

		for {
			select {
			case m := <-t.queue:
				t.write(m)

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

func (t *File) Debug(args ...interface{}) {
	if t.Level <= interfaces.DEBUG {
		t.queue <- fmt.Sprint(args...)

		if t.Parent != nil {
			t.Parent.Debug(args)
		}
	}
}

func (t *File) Info(args ...interface{}) {
	if t.Level <= interfaces.INFO {
		t.queue <- fmt.Sprint(args...)

		if t.Parent != nil {
			t.Parent.Info(args)
		}
	}
}

func (t *File) Warn(args ...interface{}) {
	if t.Level <= interfaces.WARNING {
		t.queue <- fmt.Sprint(args...)

		if t.Parent != nil {
			t.Parent.Warn(args)
		}
	}
}

func (t *File) Message(args ...interface{}) {
	if t.Level <= interfaces.MESSAGE {
		t.queue <- fmt.Sprint(args...)

		if t.Parent != nil {
			t.Parent.Message(args)
		}
	}
}

func (t *File) Error(err error) {
	if t.Level <= interfaces.ERROR {
		t.queue <- fmt.Sprint(err)

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *File) Fatal(err error) {
	if t.Level <= interfaces.CRITICAL {
		t.queue <- fmt.Sprint(err)

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *File) Panic(err error) {
	if t.Level <= interfaces.CRITICAL {
		t.queue <- fmt.Sprint(err)

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
