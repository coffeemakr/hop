package handlers

import (
	"context"
	"fmt"
	"github.com/coffeemakr/wedo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type taskWithGroupModel struct {
	Groups []*wedo.Group `bson:"groups"`
	Task wedo.Task
}

var moveToTasks = bson.D{
	{    "$replaceWith", bson.D{{"task", "$$ROOT"}} },
}

var lookupGroupForTask = bson.D{
	{"$lookup", bson.D{
		{"from", "groups"},
		{"localField", "task.groupid"},
		{"foreignField", "id"},
		{"as", "groups"},
	}},
}

func createTask(ctx context.Context, task *wedo.Task) (err error) {
	task.ID = generateId()
	var taskToStore = *task
	taskToStore.Group = nil
	taskToStore.Assignee = nil
	_, err = taskCollection.InsertOne(ctx, taskToStore)
	return
}

func updateTaskById(ctx context.Context, task *wedo.Task) error {
	updateResult, err := taskCollection.UpdateOne(ctx, bson.D{{"id", task.ID}}, task)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return ErrNoSuchTask
	}
	return nil
}

func createTaskExecution(ctx context.Context, execution *wedo.TaskExecution) error {
	_, err := taskExecutionCollection.InsertOne(ctx, execution)
	return err
}

func getTaskByIdIncludingGroup(ctx context.Context, taskId string) (*wedo.Task, error) {

	match := bson.D{{"$match", bson.D{{"id", taskId},}}}
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cursor, err := taskCollection.Aggregate(ctx, mongo.Pipeline{match, moveToTasks, lookupGroupForTask}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if !cursor.Next(ctx) {
		return nil, ErrNoSuchTask
	}
	task, err := decodeTaskWithGroup(cursor)
	if err != nil {
		return nil, err
	}
	if cursor.Next(ctx) {
		return nil, ErrMultipleTaskedMatched
	}
	log.Println("task", task)
	return task, nil
}

func decodeTaskWithGroup(cursor *mongo.Cursor) (*wedo.Task, error) {

	var task *wedo.Task
	var dbTask taskWithGroupModel
	err := cursor.Decode(&dbTask)
	if err != nil {
		return nil, fmt.Errorf("decoding task failed: %s", err)
	}
	task = &dbTask.Task
	task.Group = dbTask.Groups[0]
	return task, nil
}

func getTasksForUser(ctx context.Context, userName string) (result []*wedo.Task, err error) {
	match := bson.D{{"$match", bson.D{
		{"groups." + memberNamesField, bson.D{
			{"$in", []string{userName}},
		}},
	}}}

	cursor, err := taskCollection.Aggregate(ctx, mongo.Pipeline{moveToTasks, lookupGroupForTask, match})
	if err != nil {
		return
	}
	for cursor.Next(ctx) {
		task, err := decodeTaskWithGroup(cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}
