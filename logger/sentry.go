package logger

import (
	"context"
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

func (t *Sentry) Debug(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.DEBUG {
		if t.Parent != nil {
			t.Parent.Debug(ctx, args...)
		}
	}
}

func (t *Sentry) Info(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.INFO {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelInfo)
			hub.CaptureMessage(fmt.Sprint(args...))
		}

		if t.Parent != nil {
			t.Parent.Info(ctx, args...)
		}
	}
}

func (t *Sentry) Warn(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.WARNING {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelWarning)
			hub.CaptureMessage(fmt.Sprint(args...))
		}

		if t.Parent != nil {
			t.Parent.Warn(ctx, args...)
		}
	}
}

func (t *Sentry) Message(ctx context.Context, args ...interface{}) {
	if t.Config.Level <= interfaces.MESSAGE {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelInfo)
			hub.CaptureMessage(fmt.Sprint(args...))
		}

		if t.Parent != nil {
			t.Parent.Message(ctx, args...)
		}
	}
}

func (t *Sentry) Error(ctx context.Context, err error) {
	if t.Config.Level <= interfaces.ERROR {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelError)
			hub.CaptureException(err)
		}

		if t.Parent != nil {
			t.Parent.Error(ctx, err)
		}
	}
}

func (t *Sentry) Fatal(ctx context.Context, err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelFatal)
			hub.CaptureException(err)
		}

		if t.Parent != nil {
			t.Parent.Fatal(ctx, err)
		}
	}
}

func (t *Sentry) Panic(ctx context.Context, err error) {
	if t.Config.Level <= interfaces.CRITICAL {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelFatal)
			hub.CaptureException(err)
		}

		if t.Parent != nil {
			t.Parent.Panic(ctx, err)
		}
	}
}

func (t *Sentry) getHub(ctx context.Context, level sentry.Level) *sentry.Hub {
	hub := sentry.CurrentHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)

		span := ctx.Value("span")
		if span != nil {
			if _, ok := span.(*sentry.Span); ok {
				scope.SetSpan(span.(*sentry.Span))
			}
		}

		tags := ctx.Value("tags")
		if tags != nil {
			for k, v := range tags.(map[string]string) {
				scope.SetTag(k, v)
			}
		}
	})

	return hub
}
