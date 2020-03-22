package wedo

import (
	"fmt"
	"time"
)

type IntervalUnit string

var (
	Never  = IntervalUnit("") // Never should not be used
	Days   = IntervalUnit("days")
	Weeks  = IntervalUnit("weeks")
	Months = IntervalUnit("monts")
	Years  = IntervalUnit("years")
)

type TaskExecution struct {
	ID         string
	ExecutorId int64
	Executor   *User
	Time       time.Time
}

type Task struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	IntervalUnit    IntervalUnit   `json:"interval_unit"`
	Interval        uint32         `json:"interval"`
	LastExecution   *TaskExecution `json:"last_execution,omitempty"`
}

func NewWeeklyTask(name string) *Task {
	return NewTask(name, Weeks, 1)
}

func NewTask(name string, intervalType IntervalUnit, interval uint32) *Task {
	return &Task{
		Name:         name,
		IntervalUnit: intervalType,
		Interval:     interval,
	}
}

func (t Task) String() string {
	return fmt.Sprintf(
		"Task{ ID=%s, Name=%s, IntervalUnit=%s, Interval=%d, LastExecution=%s }",
		t.ID, t.Name, t.IntervalUnit, t.Interval, t.LastExecution)
}
