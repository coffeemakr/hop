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
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var (
	ErrNoSuchTask            = errors.New("no such task")
	ErrMultipleTaskedMatched = errors.New("multiple tasks matched")
	HttpErrTaskNotFound      = http_error.NewHttpErrorType(http.StatusNotFound, "task not found")
	lookupGroupForTask       = bson.D{
		{"$lookup", bson.D{
			{"from", "groups"},
			{"localField", "groupid"},
			{"foreignField", "id"},
			{"as", "group"},
		}},
	}
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

func CreateTaskExecution(w http.ResponseWriter, r *http.Request) {
	taskId := getTaskId(r)
	ctx := r.Context()
	task, err := getTaskByIdIncludingGroup(ctx, taskId)
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

func getTaskByIdIncludingGroup(ctx context.Context, taskId string) (*wedo.Task, error) {
	var task wedo.Task
	match := bson.D{{"$match", bson.D{{"id", taskId},}}}
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cursor, err := taskCollection.Aggregate(ctx, mongo.Pipeline{match, lookupGroupForTask}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if !cursor.Next(ctx) {
		return nil, ErrNoSuchTask
	}
	err = cursor.Decode(&task)
	if err != nil {
		return nil, fmt.Errorf("decoding task failed: %s", err)
	}
	if cursor.Next(ctx) {
		return nil, ErrMultipleTaskedMatched
	}
	log.Println("task", task)
	return &task, nil
}

func getTasksForUser(ctx context.Context, userName string) (result []*wedo.Task, err error) {
	match := bson.D{{"$match", bson.D{
		{"group." + memberNamesField, bson.D{
			{"$in", []string{userName}},
		}},
	}}}

	cursor, err := taskCollection.Aggregate(ctx, mongo.Pipeline{lookupGroupForTask, match})
	if err != nil {
		return
	}
	for cursor.Next(ctx) {
		var task wedo.Task
		err = cursor.Decode(&task)
		if err != nil {
			return
		}
		result = append(result, &task)
	}
	return result, nil
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
