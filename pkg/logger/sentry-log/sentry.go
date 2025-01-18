package sentry_log

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"gitlab.com/devpro_studio/go_utils/decode"
	"time"
)

type Sentry struct {
	name   string
	parent ILogger
	config Config
	enable bool
}

type Config struct {
	Level           LogLevel `yaml:"level"`
	SentryURL       string   `yaml:"sentry_url"`
	AppEnv          string   `yaml:"app_env"`
	SampleRate      float64  `yaml:"sample_rate"`
	TraceSampleRate float64  `yaml:"trace_sample_rate"`
	Enable          bool     `yaml:"enable"`
}

func NewSentry(name string, parent ILogger) *Sentry {
	return &Sentry{
		name:   name,
		parent: parent,
	}
}

func (t *Sentry) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.AppEnv == "" {
		t.config.AppEnv = "local"
	}

	if t.config.SentryURL != "" && t.config.Enable {
		transport := sentry.NewHTTPTransport()
		transport.Timeout = time.Second * 3

		err := sentry.Init(sentry.ClientOptions{
			Dsn:              t.config.SentryURL,
			SampleRate:       t.config.SampleRate,
			TracesSampleRate: t.config.TraceSampleRate,
			Environment:      t.config.AppEnv,
			Transport:        transport,
			EnableTracing:    t.config.TraceSampleRate > 0,
		})

		if err != nil {
			fmt.Printf("sentry.Init: %s\n", err)
		} else {
			t.enable = true
		}
	}

	return nil
}

func (t *Sentry) Stop() error {
	if t.parent != nil {
		return t.parent.Stop()
	}

	if t.enable {
		sentry.Flush(time.Second * 2)
	}

	return nil
}

func (t *Sentry) Name() string {
	return t.name
}

func (t *Sentry) Type() string {
	return "logger"
}

func (t *Sentry) SetLevel(level LogLevel) {
	t.config.Level = level

	if t.parent != nil {
		t.parent.SetLevel(level)
	}
}

func (t *Sentry) Debug(ctx context.Context, args ...interface{}) {
	if t.config.Level <= DEBUG {
		if t.parent != nil {
			t.parent.Debug(ctx, args...)
		}
	}
}

func (t *Sentry) Info(ctx context.Context, args ...interface{}) {
	if t.config.Level <= INFO {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelInfo)
			hub.CaptureMessage(fmt.Sprint(args...))
		}

		if t.parent != nil {
			t.parent.Info(ctx, args...)
		}
	}
}

func (t *Sentry) Warn(ctx context.Context, args ...interface{}) {
	if t.config.Level <= WARNING {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelWarning)
			hub.CaptureMessage(fmt.Sprint(args...))
		}

		if t.parent != nil {
			t.parent.Warn(ctx, args...)
		}
	}
}

func (t *Sentry) Message(ctx context.Context, args ...interface{}) {
	if t.config.Level <= MESSAGE {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelInfo)
			hub.CaptureMessage(fmt.Sprint(args...))
		}

		if t.parent != nil {
			t.parent.Message(ctx, args...)
		}
	}
}

func (t *Sentry) Error(ctx context.Context, err error) {
	if t.config.Level <= ERROR {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelError)
			hub.CaptureException(err)
		}

		if t.parent != nil {
			t.parent.Error(ctx, err)
		}
	}
}

func (t *Sentry) Fatal(ctx context.Context, err error) {
	if t.config.Level <= CRITICAL {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelFatal)
			hub.CaptureException(err)
		}

		if t.parent != nil {
			t.parent.Fatal(ctx, err)
		}
	}
}

func (t *Sentry) Panic(ctx context.Context, err error) {
	if t.config.Level <= CRITICAL {
		if t.enable {
			hub := t.getHub(ctx, sentry.LevelFatal)
			hub.CaptureException(err)
		}

		if t.parent != nil {
			t.parent.Panic(ctx, err)
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

func (t *Sentry) Parent() interface{} {
	return t.parent
}

func (t *Sentry) SetParent(parent interface{}) {
	t.parent = parent.(ILogger)
}
