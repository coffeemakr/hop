package cmd

import (
	"context"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	crypto_rand "crypto/rand"
	math_rand "math/rand"

	"github.com/coffeemakr/ruck/server/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/square/go-jose/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	serverHTTPPort = 8080
	serverHTTPHost = "127.0.0.1"
	databaseURL    = "mongodb://wedo:secret@localhost:27017"
	authenticator  *handlers.Authenticator

	serverCommand = &cobra.Command{
		Use:  "serve",
		RunE: runServer,
	}
)

func init() {
	serverCommand.PersistentFlags().IntVar(&serverHTTPPort, "http-port", serverHTTPPort, "The port to listen on")
	serverCommand.PersistentFlags().StringVar(&serverHTTPHost, "http-host", serverHTTPHost, "The host to listen on")
	serverCommand.PersistentFlags().StringVar(&databaseURL, "database-url", databaseURL, "The host to listen on")

	viper.BindPFlag("listen.port", serverCommand.PersistentFlags().Lookup("http-port"))
	viper.BindPFlag("listen.host", serverCommand.PersistentFlags().Lookup("http-host"))
	viper.BindPFlag("database.url", serverCommand.PersistentFlags().Lookup("database-url"))
}

func getServerAddress() string {
	addr := serverHTTPHost + ":" + strconv.FormatInt(int64(serverHTTPPort), 10)
	return addr
}

func securelySeed() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func runServer(*cobra.Command, []string) error {

	securelySeed()

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

	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURL))
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

	addr := getServerAddress()
	log.Printf("Starting server at %s\n", addr)
	router := mux.NewRouter()
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")

	api := router.MatcherFunc(func(request *http.Request, match *mux.RouteMatch) bool {
		return "" != request.Header.Get("Authorization")
	}).Subrouter()
	api.HandleFunc("/groups", handlers.GetAllGroups).Methods("GET")
	api.HandleFunc("/groups", handlers.CreateGroup).Methods("POST")
	api.HandleFunc("/groups/{groupId}", handlers.GetGroup).Methods("GET")
	api.HandleFunc("/groups/{groupId}", handlers.DeleteGroup).Methods("DELETE")
	api.HandleFunc("/groups/{groupId}/join", handlers.JoinGroup).Methods("POST")
	api.HandleFunc("/groups/{groupId}/tasks", handlers.CreateTaskForGroup).Methods("POST")
	api.HandleFunc("/tasks", handlers.GetAllTasks).Methods("GET")
	api.HandleFunc("/tasks/{taskId}", handlers.GetTaskById).Methods("GET")
	api.HandleFunc("/tasks/{taskId}/complete", handlers.CreateTaskExecution).Methods("POST")
	api.Use(authenticator.MiddleWare)

	return http.ListenAndServe(addr, router)
}

func loadPrivateKey() (*jose.JSONWebKey, error) {
	var key jose.JSONWebKey
	fp, err := os.Open(jwkName)
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
