package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/url"
	"os"
)

var rootCommand = &cobra.Command{
	Use: "wedo",
}

var (
	client *Client
	proxyStr string = "http://localhost:9090"
)
func init() {
	rootCommand.AddCommand(doneCommand, loginCommand, registerCommand, completionCommand, groupCommand,
		taskCommand)

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	client = &Client{
		BaseUrl:    "http://localhost:8080",
		Client:     &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		},
		TokenStore: NewFileTokenStore(os.ExpandEnv("$HOME/.wedo-cred.txt")),
	}
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
