package handlers

import (
	"context"
	"encoding/json"
	"errors"
	http_error "github.com/coffeemakr/go-http-error"
	"github.com/coffeemakr/amtli"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	ErrNoSuchTask             = errors.New("no such task")
	ErrMultipleTaskedMatched  = errors.New("multiple tasks matched")
	HttpErrTaskNotFound       = http_error.NewHttpErrorType(http.StatusNotFound, "task not found")
	HttpErrAssigneeNotInGroup = http_error.NewHttpErrorType(http.StatusBadRequest, "assignee not in group")
)

func getGroupId(r *http.Request) string {
	var vars = mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok || groupId == "" {
		panic("Can't read group id")
	}
	return groupId
}

func getTaskId(r *http.Request) string {
	var vars = mux.Vars(r)
	groupId, ok := vars["taskId"]
	if !ok || groupId == "" {
		panic("Can't read task id")
	}
	return groupId
}

func randomChoice(values []string) string {
	index := rand.Intn(len(values))
	return values[index]
}

func stringArrayContain(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CreateTaskForGroup(w http.ResponseWriter, r *http.Request) {
	var (
		task  amtli.Task
		group *amtli.Group
		err   error
		ctx   = r.Context()
	)
	groupId := getGroupId(r)
	userName, err := GetUserNameFromRequest(r)
	if err != nil {
		panic(err)
	}
	// decode task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		ErrInvalidJsonBody.Cause(err).Write(w, r)
		return
	}

	// load group and check therefore if the user is a member of the group
	group, err = getGroupForUser(ctx, groupId, userName)
	if err != nil {
		http_error.ErrBadRequest.Cause(err).Write(w, r)
		return
	}

	if task.AssigneeName == "" {
		task.AssigneeName = randomChoice(group.MemberNames)
	} else if !stringArrayContain(group.MemberNames, task.AssigneeName) {
		HttpErrAssigneeNotInGroup.Causef("can't assign %s", task.AssigneeName).Write(w, r)
		return
	}

	switch task.Interval.Unit {
	case amtli.Days:
	case amtli.Months:
	case amtli.Years:
	case amtli.Weeks:
	// Ok
	default:
		http_error.ErrBadRequest.CauseString("Invalid interval unit").Write(w, r)
		return

	}

	task.LastExecution = nil
	task.DueDate = task.Interval.Next(time.Now())
	task.GroupID = groupId
	if err := createTask(ctx, &task); err != nil {
		http_error.ErrInternalServerError.Causef("Failed to create task: %s", err).Write(w, r)
		return
	}
	if err := writeJson(w, task); err != nil {
		http_error.ErrInternalServerError.Causef("Failed to write response: %s", err).Write(w, r)
	}
}

func UpdateTaskById(w http.ResponseWriter, r *http.Request) {
	taskId := getTaskId(r)
	ctx := r.Context()
	var updateTask amtli.Task
	err := json.NewDecoder(r.Body).Decode(&updateTask)
	if err != nil {
		http_error.ErrBadRequest.Cause(err).Write(w, r)
		return
	}
	updateTask.ID = taskId
	err = updateTaskById(ctx, &updateTask)
	switch err {
	case ErrNoSuchTask:
		HttpErrTaskNotFound.Cause(err).Write(w, r)
	case nil:
		mustWriteJson(w, updateTask)
	default:
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
	}
}

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	taskId := getTaskId(r)
	ctx := r.Context()
	task, err := getTaskByIdIncludingGroup(ctx, taskId)
	switch err {
	case ErrNoSuchTask:
		HttpErrTaskNotFound.Cause(err).Write(w, r)
	case nil:
		mustWriteJson(w, task)
	default:
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
	}
}

func assignTaskToNextPerson(ctx context.Context, executorName string, task *amtli.Task) error {
	// TODO: The one that executed the task should be put at the end of the queue
	if task.AssigneeName == executorName {
		task.AssignNext()
	}
	task.DueDate = task.Interval.Next(time.Now())
	return updateTaskById(ctx, task)
}

func CreateTaskExecution(w http.ResponseWriter, r *http.Request) {
	taskId := getTaskId(r)
	userName, err := GetUserNameFromRequest(r)
	if err != nil {
		panic(err)
	}
	ctx := r.Context()
	task, err := getTaskByIdIncludingGroup(ctx, taskId)
	if err != nil {
		log.Printf("Failed to load get task for execution: %s", err)
		if err == ErrNoSuchTask {
			HttpErrTaskNotFound.Cause(err).Write(w, r)
		} else {
			http_error.ErrInternalServerError.Cause(err).Write(w, r)
		}
		return
	}

	execution := amtli.TaskExecution{
		ExecutorName: userName,
		Time:         time.Now(),
		TaskId:       taskId,
		Task:         task,
	}
	if err := createTaskExecution(ctx, &execution); err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
		return
	}

	if err := assignTaskToNextPerson(ctx, userName, task); err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
	}
	log.Printf("Created task execution: %s\n", execution)
	mustWriteJson(w, execution)
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUserNameFromRequest(r)
	if err != nil {
		panic(err)
	}
	tasks, err := getTasksForUser(r.Context(), userName)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
		return
	}
	mustWriteJson(w, tasks)
}
