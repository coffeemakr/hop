package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	http_error "github.com/coffeemakr/go-http-error"
	"github.com/coffeemakr/wedo"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)
var (
	ErrNoSuchTask = errors.New("no such task")
	HttpErrTaskNotFound = http_error.NewHttpErrorType(http.StatusNotFound, "task not found")

)

func createTask(ctx context.Context, task *wedo.Task) (err error) {
	task.ID = generateId()
	_, err = taskCollection.InsertOne(ctx, task)
	return
}

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

func CreateTaskForGroup(w http.ResponseWriter, r *http.Request) {
	var task wedo.Task
	groupId := getGroupId(r)
	var ctx = r.Context()
	// TODO: check if user has write permissions on group
	var _, err = GetUserNameFromRequest(r)
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		ErrInvalidJsonBody.Cause(err).Write(w, r)
		return
	}
	task.LastExecution = nil
	task.GroupID = groupId
	if err := createTask(ctx, &task); err != nil {
		http_error.ErrInternalServerError.Causef("Failed to create task: %s", err).Write(w, r)
		return
	}
	if err := writeJson(w, task); err != nil {
		http_error.ErrInternalServerError.Causef("Failed to write response: %s", err).Write(w, r)
	}
}

func CreateTaskExecution(w http.ResponseWriter, r *http.Request) {
	taskId := getTaskId(r)
	ctx := r.Context()
	task, err := getTasksById(ctx, taskId)
	if err != nil {
		if err == ErrNoSuchTask {
			HttpErrTaskNotFound.Cause(err).Write(w, r)
		} else {
			http_error.ErrInternalServerError.Cause(err).Write(w, r)
		}
		return
	}
	fmt.Fprintf(w, "Task execution! %s", task)
}

func getTasksById(ctx context.Context, taskId string) (*wedo.Task, error) {
	var task wedo.Task
	err := taskCollection.FindOne(ctx, bson.M{"id": bson.M{"$eq": taskId}}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchTask
		}
		return nil, err
	}
	return &task, nil
}
