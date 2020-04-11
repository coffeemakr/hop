package main

import (
	"bufio"
	"context"
	"crypto"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coffeemakr/amtli/handlers"
	"github.com/gorilla/mux"
	"github.com/square/go-jose/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	crypto_rand "crypto/rand"
	"crypto/rsa"
	math_rand "math/rand"

	"github.com/spf13/cobra"
)

const (
	jwkName = "private-server-jwk.json"
)

var (
	serverHttpPort = 8080
	serverHttpHost = ""
	rootCmd        = &cobra.Command{
		Use: "hop-server",
	}
	serverCommand = &cobra.Command{
		Use:  "serve",
		RunE: runServer,
	}
	authenticator       *handlers.Authenticator
	generateKeysCommand *cobra.Command = &cobra.Command{
		Use:  "generate-secrets",
		RunE: generateKeys,
	}
)

func init() {
	rootCmd.AddCommand(generateKeysCommand, serverCommand)
}

func generateSignatureRsaKey() (*jose.JSONWebKey, error) {
	key, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	priv := jose.JSONWebKey{
		Key:       key,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}

	thumb, err := priv.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	priv.KeyID = base64.RawURLEncoding.EncodeToString(thumb)
	return &priv, nil
}

func askYesOrNo(question string) (bool, error) {
	var valid, answer bool
	for !valid {
		print(question)
		print(" [y/n] ")
		reader := bufio.NewReader(os.Stdin)
		answerString, _ := reader.ReadString('\n')
		answerString = answerString[:len(answerString)-1]
		switch answerString {
		case "y", "Y":
			answer = true
			valid = true
		case "n", "N":
			answer = false
			valid = true
		default:
			fmt.Printf("\nInvalid answer. Type 'y' or 'n'.\n")
		}
	}
	return answer, nil
}

func generateKeys(*cobra.Command, []string) error {
	priv, err := generateSignatureRsaKey()
	if err != nil {
		return err
	}
	privJSON, err := priv.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = os.Stat(jwkName)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Failed to check if file exists: %s\n", err)
	}
	if err == nil {
		answer, err := askYesOrNo(fmt.Sprintf("File %s already exists.\nOverwrite?", jwkName))
		if err != nil {
			return err
		}
		if !answer {
			fmt.Printf("User abort!\n")
			os.Exit(1)
		}
	}
	err = ioutil.WriteFile(jwkName, privJSON, 0600)
	if err == nil {
		fmt.Printf("Key ID: %s\nKey successfully generated.\n", priv.KeyID)
		os.Exit(0)
	}
	return err
}

func getServerAddress() string {
	addr := serverHttpHost + ":" + strconv.FormatInt(int64(serverHttpPort), 10)
	return addr
}

func runServer(*cobra.Command, []string) error {

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

func securelySeed() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func main() {
	securelySeed()

	log.Fatal(rootCmd.Execute())
}
