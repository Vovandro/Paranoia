package interfaces

import (
	"context"
	"time"
)

type ITaskRunConfiguration interface {
}

type ITask interface {
	Init(app IEngine) error
	Stop() error
	Name() string
	Start() []ITaskRunConfiguration
	Invoke(ctx context.Context, data map[string]interface{})
}

type TaskRunAfter struct {
	ITaskRunConfiguration
	Restart <-chan time.Duration
	After   time.Duration
}

type TaskRunTime struct {
	ITaskRunConfiguration
	Restart <-chan time.Time
	To      time.Time
}
