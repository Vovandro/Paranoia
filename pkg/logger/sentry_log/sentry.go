package sentry_log

import (
	"context"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"gitlab.com/devpro_studio/go_utils/decode"
)

type Sentry struct {
	name   string
	parent ILogger
	config Config
	enable bool
	debug  bool
}

type Config struct {
	Level           LogLevel `yaml:"level"`
	SentryURL       string   `yaml:"sentry_url"`
	AppEnv          string   `yaml:"app_env"`
	SampleRate      float64  `yaml:"sample_rate"`
	TraceSampleRate float64  `yaml:"trace_sample_rate"`
	Enable          bool     `yaml:"enable"`
	Debug           bool     `yaml:"debug"`
}

func New(name string) *Sentry {
	return &Sentry{
		name: name,
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
			Debug:            t.config.Debug,
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

func (t *Sentry) SetLevel(level int) {
	t.config.Level = LogLevel(level)

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
			hub := t.getHub(ctx, sentry.LevelInfo, nil)
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
			hub := t.getHub(ctx, sentry.LevelWarning, nil)
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
			hub := t.getHub(ctx, sentry.LevelInfo, nil)
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
			hub := t.getHub(ctx, sentry.LevelError, err)
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
			hub := t.getHub(ctx, sentry.LevelFatal, err)
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
			hub := t.getHub(ctx, sentry.LevelFatal, err)
			hub.CaptureException(err)
		}

		if t.parent != nil {
			t.parent.Panic(ctx, err)
		}
	}
}

func (t *Sentry) getHub(ctx context.Context, level sentry.Level, err error) *sentry.Hub {
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

		// Add stack trace information to the scope
		if err != nil {
			scope.SetExtra("stack_trace", fmt.Sprintf("%+v", err))
		}
		// Add request information from context if available
		req := ctx.Value("request")
		if req != nil {
			if httpReq, ok := req.(map[string]interface{}); ok {
				for k, v := range httpReq {
					scope.SetExtra(k, v)
				}
			}
		}

		// Extract OLTP trace information from context
		if trace := ctx.Value("trace"); trace != nil {
			if traceData, ok := trace.(map[string]interface{}); ok {
				for k, v := range traceData {
					scope.SetExtra("trace."+k, v)
				}
			}
		}

		// Add trace ID if available
		if traceID := ctx.Value("trace_id"); traceID != nil {
			if id, ok := traceID.(string); ok {
				scope.SetTag("trace_id", id)
			}
		}

		// Add span ID if available
		if spanID := ctx.Value("span_id"); spanID != nil {
			if id, ok := spanID.(string); ok {
				scope.SetTag("span_id", id)
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
