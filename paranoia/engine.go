package paranoia

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/paranoia/config/yaml"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/Paranoia/paranoia/telemetry"
)

type Engine struct {
	name string

	starting bool

	config         interfaces.IConfig
	logger         interfaces.ILogger
	metricExporter interfaces.IMetrics
	trace          interfaces.ITrace

	task task

	pkg         map[string]map[string]interfaces.IPkg
	modules     map[string]map[string]interfaces.IModules
	middlewares map[string]interface{}
}

func New(name string, configName string) *Engine {
	t := &Engine{}

	t.starting = false
	t.name = name
	t.config = yaml.New(yaml.AutoConfig{FName: configName})

	t.pkg = make(map[string]map[string]interfaces.IPkg, 10)
	t.modules = make(map[string]map[string]interfaces.IModules, 10)
	t.middlewares = make(map[string]interface{})

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

func (t *Engine) GetLogger() interfaces.ILogger {
	return t.logger
}

func (t *Engine) GetConfig() interfaces.IConfig {
	return t.config
}

func (t *Engine) SetMetrics(c interfaces.IMetrics) {
	if t.metricExporter != nil {
		_ = t.metricExporter.Stop()
	}

	t.metricExporter = c
}

func (t *Engine) SetTrace(c interfaces.ITrace) {
	if t.trace != nil {
		_ = t.trace.Stop()
	}

	t.trace = c
}

func (t *Engine) PushPkg(c interfaces.IPkg) interfaces.IEngine {
	if c == nil {
		panic("nil package")
		return nil
	}

	name := c.Name()
	typePkg := c.Type()

	if typePkg == interfaces.PkgLogger {
		convertedLogger, ok := c.(interfaces.ILogger)

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

			l = l.Parent().(interfaces.ILogger)
		}

		l.SetParent(convertedLogger)
	} else {
		if p, ok := t.pkg[typePkg]; ok {
			if _, ok := p[name]; ok {
				panic("duplicate package: " + typePkg + " - " + name)
			}

			p[name] = c
		} else {
			t.pkg[typePkg] = make(map[string]interfaces.IPkg)
			t.pkg[typePkg][name] = c
		}
	}

	return t
}

func (t *Engine) GetPkg(typePkg string, key string) interfaces.IPkg {
	if p, ok := t.pkg[typePkg]; ok {
		if pkg, ok := p[key]; ok {
			return pkg
		}
	}

	return nil
}

func (t *Engine) PushModule(c interfaces.IModules) interfaces.IEngine {
	if c == nil {
		panic("nil package")
		return nil
	}

	name := c.Name()
	typeModule := c.Type()

	if typeModule == interfaces.ModuleMiddleware {
		convertedMiddleware, ok := c.(interfaces.IMiddleware)

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
			t.modules[typeModule] = make(map[string]interfaces.IModules)
			t.modules[typeModule][name] = c
		}
	}

	return t
}

func (t *Engine) GetModule(typeModule string, key string) interfaces.IModules {
	if typeModule == interfaces.ModuleMiddleware {
		if m, ok := t.middlewares[key]; ok {
			return m.(interfaces.IMiddleware)
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

func (t *Engine) GetTask(key string) interfaces.ITask {
	return t.task.GetTask(key)
}

func (t *Engine) PushTask(b interfaces.ITask) interfaces.IEngine {
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

	for ; l != nil; l = l.Parent().(interfaces.ILogger) {
		err = l.Init(t.config.GetConfigItem(l.Type(), l.Name()))

		if err != nil {
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to init logger %s: %w", l.Name(), err))
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
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to init metrics %s: %w", t.metricExporter.Name(), err))
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
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to init trace %s: %w", t.trace.Name(), err))
			return err
		}
	}

	for typePkg, pkg := range t.pkg {
		for name, c := range pkg {
			cfg := t.config.GetConfigItem(typePkg, name)
			if typePkg == interfaces.PkgServer {
				cfg["middlewares"] = t.middlewares
			}
			err = c.Init(cfg)

			if err != nil {
				t.logger.Fatal(context.Background(), fmt.Errorf("failed to init package %s %s: %w", typePkg, name, err))
				return err
			}
		}
	}

	for name, c := range t.middlewares {
		err = c.(interfaces.IMiddleware).Init(t, t.config.GetConfigItem(interfaces.ModuleMiddleware, name))

		if err != nil {
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to init middleware %s: %w", name, err))
			return err
		}
	}

	for typeModule, modules := range t.modules {
		for name, c := range modules {
			err = c.Init(t, t.config.GetConfigItem(typeModule, name))
			if err != nil {
				t.logger.Fatal(context.Background(), fmt.Errorf("failed to init module %s %s: %w", typeModule, name, err))
				return err
			}
		}
	}

	t.task.Start()

	if servers, ok := t.pkg[interfaces.PkgServer]; ok {
		for _, server := range servers {
			err = server.(interfaces.IServer).Start()

			if err != nil {
				t.logger.Fatal(context.Background(), fmt.Errorf("failed to start server %s: %w", server.Name(), err))
				return err
			}
		}
	}

	if t.trace != nil {
		err = t.trace.Start()

		if err != nil {
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to start trace %s: %w", t.trace.Name(), err))
			return err
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Start()

		if err != nil {
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to start metric exporter %s: %w", t.metricExporter.Name(), err))
			return err
		}
	}

	t.starting = true

	return err
}

func (t *Engine) Stop() error {
	var err error = nil

	t.starting = false

	if servers, ok := t.pkg[interfaces.PkgServer]; ok {
		for _, server := range servers {
			err = server.Stop()

			if err != nil {
				t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop server %s: %w", server.Name(), err))
				return err
			}
		}
	}

	t.task.Stop()

	for typeModule, modules := range t.modules {
		for name, m := range modules {
			err = m.Stop()

			if err != nil {
				t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop module %s %s: %w", typeModule, name, err))
				return err
			}
		}
	}

	for typePkg, pkg := range t.pkg {
		if typePkg == interfaces.PkgServer {
			continue
		}
		for name, p := range pkg {
			err = p.Stop()

			if err != nil {
				t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop package %s %s: %w", typePkg, name, err))
				return err
			}
		}
	}

	for name, item := range t.middlewares {
		err = item.(interfaces.IMiddleware).Stop()

		if err != nil {
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop middleware %s: %w", name, err))
			return err
		}
	}

	if t.metricExporter != nil {
		err = t.metricExporter.Stop()

		if err != nil {
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop metric exporter %s: %w", t.metricExporter.Name(), err))
			return err
		}
	}

	if t.trace != nil {
		err = t.trace.Stop()

		if err != nil {
			t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop trace %s: %w", t.trace.Name(), err))
			return err
		}
	}

	err = t.config.Stop()

	if err != nil {
		t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop config: %w", err))
		return err
	}

	err = t.logger.Stop()

	if err != nil {
		t.logger.Fatal(context.Background(), fmt.Errorf("failed to stop logger %s: %w", t.logger.Name(), err))
		return err
	}

	return err
}
