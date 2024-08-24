package Paranoia

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"sync/atomic"
	"testing"
	"time"
)

type testTask struct {
	count atomic.Int32
	cfg   []interfaces.ITaskRunConfiguration
}

func (t *testTask) Init(app interfaces.IService) error        { return nil }
func (t *testTask) Stop() error                               { return nil }
func (t *testTask) String() string                            { return "test" }
func (t *testTask) Start() []interfaces.ITaskRunConfiguration { return t.cfg }
func (t *testTask) Invoke(map[string]interface{})             { t.count.Add(1) }

func Test_task_run(t1 *testing.T) {
	t1.Run("base test", func(t1 *testing.T) {
		reset := make(chan time.Duration, 1)
		defer close(reset)

		tsk := &testTask{
			cfg: []interfaces.ITaskRunConfiguration{
				&interfaces.TaskRunAfter{
					Restart: reset,
					After:   time.Millisecond * 1000,
				},
			},
		}

		t := task{}
		t.Init(nil)
		t.Start()

		t.PushTask(tsk, true)

		time.Sleep(time.Second)
		reset <- time.Millisecond * 1000

		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)

			if tsk.count.Load() == 2 {
				break
			}
		}

		t.Stop()

		if tsk.count.Load() != 2 {
			t1.Errorf("expect 2 tasks, got %d", tsk.count.Load())
		}
	})
}

func Test_task_ManualRun(t1 *testing.T) {
	t1.Run("base test", func(t1 *testing.T) {
		tsk := &testTask{
			cfg: []interfaces.ITaskRunConfiguration{},
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
