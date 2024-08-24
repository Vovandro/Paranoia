package interfaces

import "time"

type ITaskRunConfiguration interface {
}

type ITask interface {
	Init(app IEngine) error
	Stop() error
	String() string
	Start() []ITaskRunConfiguration
	Invoke(map[string]interface{})
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
