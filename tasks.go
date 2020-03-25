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

type Interval struct {
	Unit   IntervalUnit `json:"unit"`
	Amount uint32       `json:"amount"`
}

func (i Interval) String() string {
	return fmt.Sprintf("every %d %s", i.Amount, i.Unit)
}

func (i Interval) Next(day time.Time) time.Time {
	switch i.Unit {
	case Days:
		return day.AddDate(0, 0, int(i.Amount))
	case Weeks:
		return day.AddDate(0, 0, int(7*i.Amount))
	case Months:
		return day.AddDate(0, int(i.Amount), 0)
	case Years:
		return day.AddDate(int(i.Amount), 0, 0)
	default:
		panic("invalid unit: " + i.Unit)
	}
}

type Task struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Interval      Interval       `json:"interval"`
	LastExecution *TaskExecution `json:"last_execution,omitempty"`
	GroupID       string         `json:"group_id,omitempty"`
	Group         *Group         `json:"group,omitempty" bson:"-"`
	Assignee      *User          `json:"assignee,omitempty" bson:"-"`
	AssigneeName  string         `json:"assignee_name,omitempty"`
	DueDate       time.Time      `json:"due_date"`
}

func (g *Group) NextName(after string) string {
	for i, memberName := range g.MemberNames {
		if memberName == after {
			i++
			i = i % len(g.MemberNames)
			return g.MemberNames[i]
		}
	}
	return g.MemberNames[0]
}

// AssignNext sets the assignee name to the next of the group.
// panics if the t.Group is null
func (t *Task) AssignNext() {
	if t.Group == nil {
		panic("Group not set")
	}
	t.AssigneeName = t.Group.NextName(t.AssigneeName)
	t.Assignee = nil
}

func NewWeeklyTask(name string) *Task {
	return NewTask(name, Weeks, 1)
}

func NewTask(name string, intervalType IntervalUnit, interval uint32) *Task {
	return &Task{
		Name: name,
		Interval: Interval{
			Unit:   intervalType,
			Amount: interval,
		},
	}
}

func (t Task) String() string {
	return fmt.Sprintf(
		"Task{ ID=%s, Name=%s, Interval=%s, LastExecution=%s DueDate=%s AssigneeName=%s groupId=%s }",
		t.ID, t.Name, t.Interval, t.LastExecution, t.DueDate, t.AssigneeName, t.GroupID)
}

func (t *Task) ShortID() string {
	return t.ID[:8]
}

type ByDueDate []*Task

func (a ByDueDate) Len() int           { return len(a) }
func (a ByDueDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDueDate) Less(i, j int) bool { return a[i].DueDate.Before(a[j].DueDate) }
