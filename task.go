package Paranoia

import (
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"sync"
	"sync/atomic"
	"time"
)

type taskRun struct {
	cfg    interfaces.ITaskRunConfiguration
	c      *time.Timer
	enable atomic.Bool
}

type task struct {
	tasks     map[string]interfaces.ITask
	runCfg    map[string][]taskRun
	taskMutex sync.RWMutex
	app       interfaces.IEngine

	done chan interface{}
	end  sync.WaitGroup
}

func (t *task) Init(app interfaces.IEngine) {
	t.app = app
	t.tasks = make(map[string]interfaces.ITask, 20)
	t.runCfg = make(map[string][]taskRun, 20)
	t.taskMutex = sync.RWMutex{}

	t.done = make(chan interface{})
}

func (t *task) GetTask(key string) interfaces.ITask {
	t.taskMutex.RLock()
	defer t.taskMutex.RUnlock()

	return t.tasks[key]
}

func (t *task) PushTask(b interfaces.ITask, run bool) {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	if item, ok := t.tasks[b.String()]; ok {
		_ = item.Stop()
	}

	t.tasks[b.String()] = b

	_ = b.Init(t.app)

	if run {
		cfgs := b.Start()

		t.runCfg[b.String()] = make([]taskRun, len(cfgs))

		for i, cfg := range cfgs {
			if c, ok := cfg.(*interfaces.TaskRunAfter); ok {
				t.runCfg[b.String()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(c.After),
					enable: atomic.Bool{},
				}

				t.runCfg[b.String()][i].enable.Store(true)
			} else if c, ok := cfg.(*interfaces.TaskRunTime); ok {
				t.runCfg[b.String()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(time.Until(c.To)),
					enable: atomic.Bool{},
				}

				t.runCfg[b.String()][i].enable.Store(true)
			}
		}
	}
}

func (t *task) RemoveTask(key string) {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	if item, ok := t.tasks[key]; ok {
		_ = item.Stop()
		delete(t.tasks, key)
	}
}

func (t *task) Start() {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	for _, item := range t.tasks {
		cfgs := item.Start()

		t.runCfg[item.String()] = make([]taskRun, len(cfgs))

		for i, cfg := range cfgs {
			if c, ok := cfg.(*interfaces.TaskRunAfter); ok {
				t.runCfg[item.String()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(c.After),
					enable: atomic.Bool{},
				}

				t.runCfg[item.String()][i].enable.Store(true)
			} else if c, ok := cfg.(*interfaces.TaskRunTime); ok {
				t.runCfg[item.String()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(time.Until(c.To)),
					enable: atomic.Bool{},
				}

				t.runCfg[item.String()][i].enable.Store(true)
			}
		}
	}

	go t.run()
}

func (t *task) Stop() {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	close(t.done)
	t.end.Wait()

	for _, item := range t.tasks {
		if _, ok := t.runCfg[item.String()]; ok {
			delete(t.runCfg, item.String())
		}

		_ = item.Stop()

		delete(t.tasks, item.String())
	}
}

func (t *task) run() {
	for {
		t.taskMutex.RLock()
		for key, configs := range t.runCfg {
			for i := 0; i < len(configs); i++ {
				if configs[i].enable.Load() {
					select {
					case <-configs[i].c.C:
						if tsk, ok := t.tasks[key]; ok {
							configs[i].enable.Store(false)
							t.end.Add(1)
							go func(tsk interfaces.ITask) {
								tsk.Invoke(nil)
								t.end.Done()
							}(tsk)
						}

					case <-t.done:
						configs[i].enable.Store(false)
						break

					default:
						break
					}
				}

				if c, ok := configs[i].cfg.(*interfaces.TaskRunAfter); ok {
					select {
					case r := <-c.Restart:
						configs[i].c.Reset(r)
						configs[i].enable.Store(true)

					default:
						break
					}
				} else if c, ok := configs[i].cfg.(*interfaces.TaskRunTime); ok {
					select {
					case r := <-c.Restart:
						configs[i].c.Reset(time.Until(r))
						configs[i].enable.Store(true)

					default:
						break
					}
				}

			}
		}
		t.taskMutex.RUnlock()

		select {
		case <-t.done:
			return

		case <-time.After(time.Millisecond * 10):
			break
		}
	}
}

func (t *task) RunTask(key string, args map[string]interface{}) error {
	t.taskMutex.RLock()
	defer t.taskMutex.RUnlock()

	if item, ok := t.tasks[key]; ok {
		t.end.Add(1)

		go func(tsk interfaces.ITask, args map[string]interface{}) {
			tsk.Invoke(args)
			t.end.Done()
		}(item, args)

		return nil
	}

	return fmt.Errorf("task not found")
}
