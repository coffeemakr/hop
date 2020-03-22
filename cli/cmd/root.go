package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var rootCommand = &cobra.Command{
	Use: "wedo",
}

var client *Client

func init() {
	rootCommand.AddCommand(doneCommand, addCommand, loginCommand, registerCommand, completionCommand)
	client = &Client{
		BaseUrl:    "http://localhost:8080",
		Client:     &http.Client{},
		TokenStore: NewFileTokenStore(os.ExpandEnv("$HOME/.wedo-cred.txt")),
	}

	err := client.LoadToken()
	switch err {
	case nil:
		log.Println("Successfully loaded token!")
	case ErrNoTokenSaved:
		log.Println("No token stored.")
	default:
		log.Fatalln(err)
	}
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
