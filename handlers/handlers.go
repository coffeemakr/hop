package handlers

import (
	"context"
	"encoding/json"
	http_error "github.com/coffeemakr/go-http-error"
	"github.com/coffeemakr/wedo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"net/http"
)

var (
	taskCollection     *mongo.Collection
	usersCollection    *mongo.Collection
	groupsCollection   *mongo.Collection
	ErrInvalidJsonBody = http_error.ErrBadRequest.WithDescription("Invalid JSON body")
)

func SetDB(db *mongo.Database) {
	taskCollection = db.Collection("tasks")
	usersCollection = db.Collection("users")
	groupsCollection = db.Collection("groups")
	_, err := usersCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: options.Index().SetName("user_name").SetUnique(true),
	})

	if err != nil {
		log.Fatal("create", err)
	}
}

func writeJson(w http.ResponseWriter, value interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
	return json.NewEncoder(w).Encode(value)
}

func getTasks(ctx context.Context) (result []*wedo.Task, err error) {
	cursor, err := taskCollection.Find(ctx, bson.D{})
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

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tasks, err := getTasks(r.Context())
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Println(err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateId() string {
	return RandStringRunes(32)
}