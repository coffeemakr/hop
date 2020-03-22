package main

import (
	"context"
	"github.com/coffeemakr/wedo/handlers"
	"github.com/gorilla/mux"
	"github.com/square/go-jose/v3"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	serverHttpPort = 8080
	serverHttpHost = ""
	rootCmd        = &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			err := runServer()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	authenticator *handlers.Authenticator
)

func getServerAddress() string {
	addr := serverHttpHost + ":" + strconv.FormatInt(int64(serverHttpPort), 10)
	return addr
}

func runServer() error {
	addr := getServerAddress()
	log.Printf("Starting server at %s\n", addr)
	router := mux.NewRouter()
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")

	api := router.MatcherFunc(func(request *http.Request, match *mux.RouteMatch) bool {
		return "" != request.Header.Get("Authorization")
	}).Subrouter()
	api.HandleFunc("/groups/{groupId}", handlers.GetGroup).Methods("GET")
	api.HandleFunc("/groups/{groupId}", handlers.DeleteGroup).Methods("DELETE")
	api.HandleFunc("/groups/{groupId}/join", handlers.JoinGroup).Methods("POST")
	api.HandleFunc("/groups", handlers.GetAllGroups).Methods("GET")
	api.HandleFunc("/groups", handlers.CreateGroup).Methods("POST")
	api.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/execution", handlers.CreateTaskExecution).Methods("POST")
	api.Use(authenticator.MiddleWare)

	return http.ListenAndServe(addr, router)
}

func loadPrivateKey() (*jose.JSONWebKey, error) {
	var key jose.JSONWebKey
	fp, err := os.Open("jwk-sig-priv.json")
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	err = key.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	key, err := loadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	handlers.UsedTokenIssuer = &handlers.JwtTokenIssuer{
		PrivateKey: key,
	}

	var keyset = &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			key.Public(),
		},
	}
	authenticator = &handlers.Authenticator{
		Verifier: &handlers.JwtTokenVerifier{KeySet: keyset},
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://wedo:secret@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("wedo")
	handlers.SetDB(db)

	log.Fatal(rootCmd.Execute())
}
