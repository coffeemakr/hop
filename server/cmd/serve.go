package cmd

import (
	"context"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	crypto_rand "crypto/rand"
	math_rand "math/rand"

	"github.com/coffeemakr/ruck/server"
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
	serverConfig   = server.Configuration{
		Listen:   &server.ListenConfig{},
		Database: &server.DatabaseConfig{},
	}
	authenticator *handlers.Authenticator
	config        *viper.Viper
	serverCommand = &cobra.Command{
		Use:     "serve",
		PreRunE: initConfig,
		RunE:    runServer,
	}
)

func initConfig(*cobra.Command, []string) error {
	err := config.ReadInConfig()
	if err != nil {
		return err
	}
	serverConfig.Listen.Port = config.GetInt("listen.port")
	serverConfig.Listen.Host = config.GetString("listen.host")
	serverConfig.Database.URL = config.GetString("database.url")
	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	serverCommand.PersistentFlags().Int("http-port", 8080, "The port to listen on")
	serverCommand.PersistentFlags().String("http-host", "127.0.0.1", "The host to listen on")
	serverCommand.PersistentFlags().String("database", "", "The host to listen on")

	config = viper.New()

	must(config.BindPFlag("listen.port", serverCommand.PersistentFlags().Lookup("http-port")))
	must(config.BindPFlag("listen.host", serverCommand.PersistentFlags().Lookup("http-host")))
	must(config.BindPFlag("database.url", serverCommand.PersistentFlags().Lookup("database")))
	config.SetConfigName("ruckd")
	config.AddConfigPath(".")
	config.AddConfigPath("/etc/ruckd")

	serverCommand.MarkFlagRequired("database")
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

	client, err := mongo.NewClient(options.Client().ApplyURI(serverConfig.Database.URL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("ruck")
	handlers.SetDB(db)

	addr := serverConfig.Listen.GetServerAddress()
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
