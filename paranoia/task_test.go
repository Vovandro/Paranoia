package paranoia

import (
	"context"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"sync/atomic"
	"testing"
	"time"
)

type testTask struct {
	count atomic.Int32
	cfg   []interfaces2.ITaskRunConfiguration
}

func (t *testTask) Init(app interfaces2.IEngine) error                      { return nil }
func (t *testTask) Stop() error                                             { return nil }
func (t *testTask) Name() string                                            { return "test" }
func (t *testTask) Start() []interfaces2.ITaskRunConfiguration              { return t.cfg }
func (t *testTask) Invoke(ctx context.Context, data map[string]interface{}) { t.count.Add(1) }

func Test_task_run(t1 *testing.T) {
	t1.Run("base test", func(t1 *testing.T) {
		reset := make(chan time.Duration, 1)
		defer close(reset)

		tsk := &testTask{
			cfg: []interfaces2.ITaskRunConfiguration{
				&interfaces2.TaskRunAfter{
					Restart: reset,
					After:   time.Millisecond * 100,
				},
			},
		}

		t := task{}
		t.Init(nil)
		t.Start()

		t.PushTask(tsk, true)
		var c int32

		for i := 0; i < 10; i++ {
			//time.Sleep(time.Second)
			<-time.After(time.Second)

			if c = tsk.count.Load(); c == 1 {
				break
			}
		}

		reset <- time.Millisecond * 100

		for i := 0; i < 10; i++ {
			//time.Sleep(time.Second)
			<-time.After(time.Second)

			if c = tsk.count.Load(); c == 2 {
				break
			}
		}

		t.Stop()

		if c != 2 {
			t1.Errorf("expect 2 tasks, got %d", c)
		}
	})
}

func Test_task_ManualRun(t1 *testing.T) {
	t1.Run("base test", func(t1 *testing.T) {
		tsk := &testTask{
			cfg: []interfaces2.ITaskRunConfiguration{},
		}

		t := task{}
		t.Init(nil)
		t.Start()

		t.PushTask(tsk, true)

		time.Sleep(time.Millisecond * 500)
		_ = t.RunTask("test", nil)
		time.Sleep(time.Second * 2)

		if tsk.count.Load() != 1 {
			t1.Errorf("expect 1 tasks, got %d", tsk.count.Load())
		}

		t.Stop()
	})
}
