package paranoia

import (
	"context"
	"fmt"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"sync"
	"sync/atomic"
	"time"
)

type taskRun struct {
	cfg    interfaces2.ITaskRunConfiguration
	c      *time.Timer
	enable atomic.Bool
}

type task struct {
	tasks     map[string]interfaces2.ITask
	runCfg    map[string][]taskRun
	taskMutex sync.RWMutex
	app       interfaces2.IEngine

	done chan interface{}
	end  sync.WaitGroup
}

func (t *task) Init(app interfaces2.IEngine) {
	t.app = app
	t.tasks = make(map[string]interfaces2.ITask, 20)
	t.runCfg = make(map[string][]taskRun, 20)
	t.taskMutex = sync.RWMutex{}

	t.done = make(chan interface{})
}

func (t *task) GetTask(key string) interfaces2.ITask {
	t.taskMutex.RLock()
	defer t.taskMutex.RUnlock()

	return t.tasks[key]
}

func (t *task) PushTask(b interfaces2.ITask, run bool) {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	if item, ok := t.tasks[b.Name()]; ok {
		_ = item.Stop()
	}

	t.tasks[b.Name()] = b

	_ = b.Init(t.app)

	if run {
		cfgs := b.Start()

		t.runCfg[b.Name()] = make([]taskRun, len(cfgs))

		for i, cfg := range cfgs {
			if c, ok := cfg.(*interfaces2.TaskRunAfter); ok {
				t.runCfg[b.Name()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(c.After),
					enable: atomic.Bool{},
				}

				t.runCfg[b.Name()][i].enable.Store(true)
			} else if c, ok := cfg.(*interfaces2.TaskRunTime); ok {
				t.runCfg[b.Name()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(time.Until(c.To)),
					enable: atomic.Bool{},
				}

				t.runCfg[b.Name()][i].enable.Store(true)
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

		t.runCfg[item.Name()] = make([]taskRun, len(cfgs))

		for i, cfg := range cfgs {
			if c, ok := cfg.(*interfaces2.TaskRunAfter); ok {
				t.runCfg[item.Name()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(c.After),
					enable: atomic.Bool{},
				}

				t.runCfg[item.Name()][i].enable.Store(true)
			} else if c, ok := cfg.(*interfaces2.TaskRunTime); ok {
				t.runCfg[item.Name()][i] = taskRun{
					cfg:    c,
					c:      time.NewTimer(time.Until(c.To)),
					enable: atomic.Bool{},
				}

				t.runCfg[item.Name()][i].enable.Store(true)
			}
		}
	}

	go t.run()
}

func (t *task) Stop() {
	close(t.done)
	t.end.Wait()

	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()

	for _, item := range t.tasks {
		if _, ok := t.runCfg[item.Name()]; ok {
			delete(t.runCfg, item.Name())
		}

		_ = item.Stop()

		delete(t.tasks, item.Name())
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
							go func(tsk interfaces2.ITask) {
								defer t.end.Done()
								tr := otel.Tracer("task")
								ctx, span := tr.Start(context.Background(), tsk.Name())
								defer span.End()

								tsk.Invoke(ctx, nil)
							}(tsk)
						}

					case <-t.done:
						configs[i].enable.Store(false)
						break

					default:
						break
					}
				}

				if c, ok := configs[i].cfg.(*interfaces2.TaskRunAfter); ok {
					select {
					case r := <-c.Restart:
						configs[i].c.Reset(r)
						configs[i].enable.Store(true)

					default:
						break
					}
				} else if c, ok := configs[i].cfg.(*interfaces2.TaskRunTime); ok {
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

		go func(tsk interfaces2.ITask, args map[string]interface{}) {
			defer t.end.Done()
			tr := otel.Tracer("task")
			ctx, span := tr.Start(context.Background(), tsk.Name())
			defer span.End()

			tsk.Invoke(ctx, args)
		}(item, args)

		return nil
	}

	return fmt.Errorf("task not found")
}
