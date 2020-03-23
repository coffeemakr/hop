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
	"log"
	"net/http"
)

var (
	HttpErrGroupNotFound = http_error.NewHttpErrorType(http.StatusNotFound, "Group not found")
	ErrGroupNotFound     = errors.New("group not found")
)

const memberNamesField = "membernames"

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var userName, err = GetUserNameFromRequest(r)
	if err != nil {
		panic(err)
	}
	var group wedo.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		ErrInvalidJsonBody.Cause(err).Write(w, r)
		return
	}
	group.MemberNames = []string{userName}
	group.ID = generateId()
	err = createGroup(ctx, &group)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
		return
	}
	err = writeJson(w, group)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
	}
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var userName, err = GetUserNameFromRequest(r)
	if err != nil {
		panic(err)
	}

	groups, err := getGroups(ctx, userName)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
		return
	}
	if err := writeJson(w, groups); err != nil {
		http_error.ErrInternalServerError.Causef("Can't write group: %s", err).Write(w, r)
		return
	}
}

func GetGroup(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var userName, err = GetUserNameFromRequest(r)
	if err != nil {
		panic(err)
	}
	var vars = mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		panic("Can't read group id")
	}
	group, err := getGroupForUser(ctx, groupId, userName)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
		return
	}
	if err := writeJson(w, group); err != nil {
		http_error.ErrInternalServerError.Causef("Can't write group: %s", err).Write(w, r)
		return
	}
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var userName, err = GetUserNameFromRequest(r)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
	}
	var vars = mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		panic("Can't read group id")
	}
	err = deleteGroupForUser(ctx, groupId, userName)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func JoinGroup(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var userName, err = GetUserNameFromRequest(r)
	if err != nil {
		http_error.ErrInternalServerError.Cause(err).Write(w, r)
	}
	var vars = mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		panic("Can't read group id")
	}
	err = joinGroup(ctx, userName, groupId)
	if err != nil {
		HttpErrGroupNotFound.Cause(err).Write(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func joinGroup(ctx context.Context, userName string, groupId string) error {
	result, err := groupsCollection.UpdateOne(ctx, bson.M{
		"id": bson.M{"$eq": groupId},
	}, bson.M{
		"$addToSet": bson.M{memberNamesField: userName},
	})
	if err != nil {
		return fmt.Errorf("joining group failed: %s", err)
	}
	if result.MatchedCount == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func deleteGroupForUser(ctx context.Context, groupId string, userName string) error {
	_, err := groupsCollection.DeleteOne(ctx, bson.M{
		"id":             bson.D{{"$eq", groupId}},
		memberNamesField: bson.D{{"$in", []string{userName}}},
	})
	return err
}

func createGroup(ctx context.Context, group *wedo.Group) error {
	_, err := groupsCollection.InsertOne(ctx, group)
	if err != nil {
		return err
	}
	return nil
}

func getGroupForUser(ctx context.Context, groupId string, userName string) (*wedo.Group, error) {
	var group wedo.Group
	err := groupsCollection.FindOne(ctx, bson.M{
		"id": groupId,
		memberNamesField: bson.M{"$in": []string{userName}},
	}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrGroupNotFound
		}
		return nil, err
	}

	return &group, nil
}

func getGroups(ctx context.Context, userName string) ([]*wedo.Group, error) {
	var results []*wedo.Group
	cursor, err := groupsCollection.Find(ctx, bson.M{
		memberNamesField: bson.M{"$in": []string{userName}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		log.Println(cursor.Current)
		var group wedo.Group
		err := cursor.Decode(&group)
		if err != nil {
			return nil, err
		}
		results = append(results, &group)
	}
	return results, nil
}
