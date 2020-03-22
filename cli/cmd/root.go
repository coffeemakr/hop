package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var rootCommand = &cobra.Command{
	Use:"wedo",
}

var client *Client

func init() {
	rootCommand.AddCommand(doneCommand, addCommand, loginCommand, registerCommand)
	client = &Client{
		BaseUrl: "http://localhost:8080",
		Client:  &http.Client{},
	}
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}