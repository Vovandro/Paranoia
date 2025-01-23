package paranoia

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/paranoia/config/yaml"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/paranoia/telemetry"
)

type Engine struct {
	name string

	starting bool

	config         interfaces2.IConfig
	logger         interfaces2.ILogger
	metricExporter interfaces2.IMetrics
	trace          interfaces2.ITrace

	task task

	pkg         map[string]map[string]interfaces2.IPkg
	modules     map[string]map[string]interfaces2.IModules
	middlewares map[string]interfaces2.IMiddleware
}

func New(name string, configName string) *Engine {
	t := &Engine{}

	t.starting = false
	t.name = name
	t.config = yaml.New(yaml.AutoConfig{FName: configName})

	t.pkg = make(map[string]map[string]interfaces2.IPkg, 10)
	t.modules = make(map[string]map[string]interfaces2.IModules, 10)
	t.middlewares = make(map[string]interfaces2.IMiddleware)

	t.task.Init(t)

	if t.config != nil {
		err := t.config.Init(t)

		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	return t
}

func (t *Engine) GetLogger() interfaces2.ILogger {
	return t.logger
}

func (t *Engine) GetConfig() interfaces2.IConfig {
	return t.config
}

func (t *Engine) SetMetrics(c interfaces2.IMetrics) {
	if t.metricExporter != nil {
		_ = t.metricExporter.Stop()
	}

	t.metricExporter = c
}

func (t *Engine) SetTrace(c interfaces2.ITrace) {
	if t.trace != nil {
		_ = t.trace.Stop()
	}

	t.trace = c
}

func (t *Engine) PushPkg(c interfaces2.IPkg) interfaces2.IEngine {
	if c == nil {
		panic("nil package")
		return nil
	}

	name := c.Name()
	typePkg := c.Type()

	if typePkg == interfaces2.PkgLogger {
		convertedLogger, ok := c.(interfaces2.ILogger)

		if !ok {
			panic("cannot cast logger")
		}

		l := t.logger

		if l == nil {
			t.logger = convertedLogger
			return t
		}

		for {
			if l.Parent() == nil {
				break
			}

			l = l.Parent().(interfaces2.ILogger)
		}

		l.SetParent(convertedLogger)
	} else {
		if p, ok := t.pkg[typePkg]; ok {
			if _, ok := p[name]; ok {
				panic("duplicate package: " + typePkg + " - " + name)
			}

			p[name] = c
		} else {
			t.pkg[typePkg] = make(map[string]interfaces2.IPkg)
			t.pkg[typePkg][name] = c
		}
	}

	return t
}

func (t *Engine) GetPkg(typePkg string, key string) interfaces2.IPkg {
	if p, ok := t.pkg[typePkg]; ok {
		if pkg, ok := p[key]; ok {
			return pkg
		}
	}

	return nil
}

func (t *Engine) PushModule(c interfaces2.IModules) interfaces2.IEngine {
	if c == nil {
		panic("nil package")
		return nil
	}

	name := c.Name()
	typeModule := c.Type()

	if typeModule == interfaces2.ModuleMiddleware {
		convertedMiddleware, ok := c.(interfaces2.IMiddleware)

		if !ok {
			panic("cannot cast middleware")
		}

		if _, ok := t.middlewares[name]; ok {
			panic("duplicate middleware: " + name)
		}

		t.middlewares[name] = convertedMiddleware
	} else {
		if p, ok := t.modules[typeModule]; ok {
			if _, ok := p[name]; ok {
				panic("duplicate modules: " + typeModule + " - " + name)
			}

			p[name] = c
		} else {
			t.modules[typeModule] = make(map[string]interfaces2.IModules)
			t.modules[typeModule][name] = c
		}
	}

	return t
}

func (t *Engine) GetModule(typeModule string, key string) interfaces2.IModules {
	if typeModule == interfaces2.ModuleMiddleware {
		if m, ok := t.middlewares[key]; ok {
			return m
		}
	} else {
		if p, ok := t.modules[typeModule]; ok {
			if m, ok := p[key]; ok {
				return m
			}
		}
	}

	return nil
}

func (t *Engine) GetTask(key string) interfaces2.ITask {
	return t.task.GetTask(key)
}

func (t *Engine) PushTask(b interfaces2.ITask) interfaces2.IEngine {
	t.task.PushTask(b, t.starting)

	return t
}

func (t *Engine) RemoveTask(key string) {
	t.task.RemoveTask(key)
}

func (t *Engine) RunTask(key string, args map[string]interface{}) error {
	return t.task.RunTask(key, args)
}

func (t *Engine) Init() error {
	var err error = nil

	l := t.logger

	for ; l != nil; l = l.Parent().(interfaces2.ILogger) {
		err = l.Init(t.config.GetConfigItem(l.Type(), l.Name()))

		if err != nil {
			t.logger.Fatal(context.Background(), err)
			return err
		}

		if l.Parent() == nil {
			break
		}
	}

	if t.metricExporter == nil {
		cfg := t.config.GetConfigItem("metrics", "")

		if len(cfg) > 0 {
			t.metricExporter = telemetry.NewMetrics(cfg)
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Init(t.config.GetConfigItem("metrics", t.metricExporter.Name()))

		if err != nil {
			t.logger.Fatal(context.Background(), err)
			return err
		}
	}

	if t.trace == nil {
		cfg := t.config.GetConfigItem("trace", "")

		if len(cfg) > 0 {
			t.trace = telemetry.NewTrace(cfg)
		}
	}

	if t.trace != nil {
		err = t.trace.Init(t.config.GetConfigItem("trace", t.trace.Name()))

		if err != nil {
			t.logger.Fatal(context.Background(), err)
			return err
		}
	}

	for typePkg, pkg := range t.pkg {
		for name, c := range pkg {
			err = c.Init(t.config.GetConfigItem(typePkg, name))

			if err != nil {
				t.logger.Fatal(context.Background(), err)
				return err
			}
		}
	}

	for name, c := range t.middlewares {
		err = c.Init(t, t.config.GetConfigItem(interfaces2.ModuleMiddleware, name))

		if err != nil {
			t.logger.Fatal(context.Background(), err)
			return err
		}
	}

	for typeModule, modules := range t.modules {
		for name, c := range modules {
			err = c.Init(t, t.config.GetConfigItem(typeModule, name))
			if err != nil {
				t.logger.Fatal(context.Background(), err)
				return err
			}
		}
	}

	t.task.Start()

	if servers, ok := t.pkg["servers"]; ok {
		for _, server := range servers {
			err = server.(interfaces2.IServer).Start()

			if err != nil {
				t.logger.Fatal(context.Background(), err)
				return err
			}
		}
	}

	if t.trace != nil {
		err = t.trace.Start()

		if err != nil {
			return err
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Start()

		if err != nil {
			return err
		}
	}

	t.starting = true

	return err
}

func (t *Engine) Stop() error {
	var err error = nil

	t.starting = false

	if servers, ok := t.pkg["servers"]; ok {
		for _, server := range servers {
			err = server.Stop()

			if err != nil {
				t.logger.Fatal(context.Background(), err)
				return err
			}
		}
	}

	t.task.Stop()

	for _, modules := range t.modules {
		for _, m := range modules {
			err = m.Stop()

			if err != nil {
				t.logger.Fatal(context.Background(), err)
				return err
			}
		}
	}

	for _, pkg := range t.pkg {
		for _, p := range pkg {
			err = p.Stop()

			if err != nil {
				t.logger.Fatal(context.Background(), err)
				return err
			}
		}
	}

	for _, item := range t.middlewares {
		err = item.Stop()

		if err != nil {
			t.logger.Fatal(context.Background(), err)
			return err
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Stop()

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if t.trace != nil {
		err = t.trace.Stop()

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	err = t.config.Stop()

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = t.logger.Stop()

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
