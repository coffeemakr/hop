package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	http_error "github.com/coffeemakr/go-http-error"
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
	// prevent browsers from displaying the JSON as HTML
	w.Header().Set("Content-Type", "application/json")
	// disable loading of any sources
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
	// disable content type sniffing, for CORB and disallow usage in script or style tags
	w.Header().Set("X-Content-Type-Options", "nosniff")
	err = json.NewEncoder(w).Encode(value)
	if err != nil {
		err = fmt.Errorf("failed to write response: %s", err)
	}
	return err
}

func mustWriteJson(w http.ResponseWriter, value interface{}) {
	err := writeJson(w, value)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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