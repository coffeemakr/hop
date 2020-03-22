package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	http_error "github.com/coffeemakr/go-http-error"
	"github.com/coffeemakr/wedo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var (
	bcryptCost                = bcrypt.DefaultCost
	HttpErrPasswordsDontMatch = http_error.ErrBadRequest.WithDescription("Passwords don't match")
	HttpErrInvalidCredentials = http_error.NewHttpErrorType(http.StatusUnauthorized, "Invalid credentials")
	ErrNoSucUser              = errors.New("No such user")
)

const userFieldName = "name"

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user *wedo.User
	var credentials wedo.Credentials
	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		ErrInvalidJsonBody.Cause(err).Write(w, r)
		return
	}

	// TODO: password policy

	user, err := getUserForCredentials(ctx, &credentials)
	if err != nil {
		if err == ErrNoSucUser {
			HttpErrInvalidCredentials.Cause(err).Write(w, r)
		} else {
			http_error.ErrInternalServerError.Causef("Failed to get user: %s", err).Write(w, r)
		}
		return
	}

	if err := writeJson(w, user); err != nil {
		http_error.ErrInternalServerError.Causef("Failed to write user response %s", err).Write(w, r)
	}
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user *wedo.User
	var registrationRequest wedo.RegistrationRequest
	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&registrationRequest); err != nil {
		ErrInvalidJsonBody.Cause(err).Write(w, r)
		return
	}

	if !bytes.Equal(registrationRequest.Password, registrationRequest.PasswordConfirmation) {
		HttpErrPasswordsDontMatch.CauseString("Password comparasion failed").Write(w, r)
		return
	}
	// TODO: password policy

	user, err := registerUser(ctx, &registrationRequest)
	if err != nil {
		http_error.ErrInternalServerError.Causef("Failed to register user: %s", err).Write(w, r)
		return
	}

	if err := writeJson(w, user); err != nil {
		http_error.ErrInternalServerError.Causef("Failed to write user response %s", err).Write(w, r)
	}
}

func createUser(ctx context.Context, user *wedo.User)  (err error) {
	_, err = usersCollection.InsertOne(ctx, user)
	user.PasswordHash = nil // Prevent hash from leaking
	return
}

func getUserWithPasswordForName(ctx context.Context, name string) (*wedo.User, error) {
	var user wedo.User
	err := usersCollection.FindOne(ctx, bson.M{userFieldName: name}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrNoSucUser
		}
		return nil, err
	}
	return &user, nil
}

func getUserForName(ctx context.Context, name string) (user *wedo.User, err error) {
	user, err = getUserWithPasswordForName(ctx, name)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = nil
	return user, err
}

func getUserForCredentials(ctx context.Context, credentials *wedo.Credentials) (user *wedo.User, err error) {
	user, err = getUserWithPasswordForName(ctx, credentials.Name)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, credentials.Password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = ErrNoSucUser
		}
		user = nil
	}
	return
}

func registerUser(ctx context.Context, registration *wedo.RegistrationRequest) (user *wedo.User, err error) {
	hashed, err := bcrypt.GenerateFromPassword(registration.Password, bcryptCost)
	if err != nil {
		return
	}
	user = &wedo.User{
		Name:          registration.Name,
		EmailAddress:  registration.Email,
		EmailVerified: false,
		PasswordHash:  hashed,
	}
	err = createUser(ctx, user)
	if err == nil {
		log.Printf("Created user %s\n", user.Name)
	}
	return
}
