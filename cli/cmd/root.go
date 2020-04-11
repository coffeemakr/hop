package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCommand = &cobra.Command{
	Use: "ruck",
}

var (
	client   *Client
	proxyStr = "http://localhost:9090"
	baseUrl  = "http://localhost:8080"
)

func init() {
	rootCommand.AddCommand(loginCommand, registerCommand, completionCommand, groupCommand, taskCommand)
	rootCommand.PersistentFlags().StringVar(&proxyStr, "proxy", "", "Proxy URL (e.g. http://localhost:8080)")
}

func Execute() {
	viper.SetConfigName("ruck-config")
	viper.AddConfigPath("$HOME/.config")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	var transport http.RoundTripper = http.DefaultClient.Transport

	if proxyStr != "" {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Println(err)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	client = &Client{
		BaseUrl: baseUrl,
		Client: &http.Client{
			Transport: transport,
		},
		TokenStore: NewFileTokenStore(os.ExpandEnv("$HOME/.ruck-cred.txt")),
	}

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
