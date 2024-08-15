package logger

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"time"
)

type Sentry struct {
	Parent interfaces.ILogger
	Config SentryConfig
	enable bool
}

type SentryConfig struct {
	Level           interfaces.LogLevel `yaml:"level"`
	SentryURL       string              `yaml:"sentry_url"`
	AppEnv          string              `yaml:"app_env"`
	SampleRate      float64             `yaml:"sample_rate"`
	TraceSampleRate float64             `yaml:"trace_sample_rate"`
}

func NewSentry(cfg SentryConfig, parent interfaces.ILogger) *Sentry {
	return &Sentry{
		Config: cfg,
		Parent: parent,
	}
}

func (t *Sentry) Init(cfg interfaces.IConfig) error {
	if cfg != nil {
		l := cfg.GetString("LOG_LEVEL", "")

		if l != "" {
			t.Config.Level = interfaces.GetLogLevel(l)
		}

		url := cfg.GetString("SENTRY_URL", "")

		if l != "" {
			t.Config.SentryURL = url
		}

		env := cfg.GetString("APP_ENV", "")

		if env != "" {
			t.Config.AppEnv = env
		}

		sr := cfg.GetFloat("SENTRY_SAMPLE_RATE", -1)

		if sr >= 0 {
			t.Config.SampleRate = sr
		}

		sr = cfg.GetFloat("SENTRY_TRACE_SAMPLE_RATE", -1)

		if sr >= 0 {
			t.Config.TraceSampleRate = sr
		}
	}

	if t.Config.AppEnv == "" {
		t.Config.AppEnv = "local"
	}

	if t.Config.SentryURL != "" {
		transport := sentry.NewHTTPTransport()
		transport.Timeout = time.Second * 3

		err := sentry.Init(sentry.ClientOptions{
			Dsn:                t.Config.SentryURL,
			SampleRate:         t.Config.SampleRate,
			TracesSampleRate:   t.Config.TraceSampleRate,
			ProfilesSampleRate: t.Config.TraceSampleRate,
			Environment:        t.Config.AppEnv,
			Transport:          transport,
			EnableTracing:      t.Config.TraceSampleRate > 0,
		})

		if err != nil {
			fmt.Printf("sentry.Init: %s\n", err)
		} else {
			t.enable = true
		}
	}

	if t.Parent != nil {
		return t.Parent.Init(cfg)
	}

	return nil
}

func (t *Sentry) Stop() error {
	if t.Parent != nil {
		return t.Parent.Stop()
	}

	if t.enable {
		sentry.Flush(time.Second * 2)
	}

	return nil
}

func (t *Sentry) SetLevel(level interfaces.LogLevel) {
	t.Config.Level = level

	if t.Parent != nil {
		t.Parent.SetLevel(level)
	}
}

func (t *Sentry) Debug(args ...interface{}) {
	if t.Config.Level <= interfaces.DEBUG {
		if t.Parent != nil {
			t.Parent.Debug(args...)
		}
	}
}

func (t *Sentry) Info(args ...interface{}) {
	if t.Config.Level <= interfaces.INFO {
		if t.Parent != nil {
			t.Parent.Info(args...)
		}
	}
}

func (t *Sentry) Warn(args ...interface{}) {
	if t.Config.Level <= interfaces.WARNING {

		if t.enable {
			hub := sentry.CurrentHub()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelWarning)
			})
			hub.CaptureMessage(fmt.Sprint(args...))
		}

		if t.Parent != nil {
			t.Parent.Warn(args...)
		}
	}
}

func (t *Sentry) Message(args ...interface{}) {
	if t.Config.Level <= interfaces.MESSAGE {
		if t.Parent != nil {
			t.Parent.Message(args...)
		}
	}
}

func (t *Sentry) Error(err error) {
	if t.Config.Level <= interfaces.ERROR {
		if t.enable {
			hub := sentry.CurrentHub()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelError)
			})
			hub.CaptureException(err)
		}

		if t.Parent != nil {
			t.Parent.Error(err)
		}
	}
}

func (t *Sentry) Fatal(err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		if t.enable {
			hub := sentry.CurrentHub()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelFatal)
			})
			hub.CaptureException(err)
		}

		if t.Parent != nil {
			t.Parent.Fatal(err)
		}
	}
}

func (t *Sentry) Panic(err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		if t.enable {
			hub := sentry.CurrentHub()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelFatal)
			})
			hub.CaptureException(err)
		}

		if t.Parent != nil {
			t.Parent.Panic(err)
		}
	}
}
