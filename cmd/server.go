package main

import (
	"context"
	"github.com/coffeemakr/wedo/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	router.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	router.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/execution", handlers.CreateTaskExecution).Methods("POST")

	return http.ListenAndServe(addr, router)
}

func main() {
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
