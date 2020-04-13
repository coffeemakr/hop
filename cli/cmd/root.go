package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/coffeemakr/ruck/cli"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use: "ruck",
}

var (
	client   *cli.Client
	proxyStr string
)

func init() {
	rootCommand.AddCommand(loginCommand, registerCommand, completionCommand, groupCommand, taskCommand, configCommand)
	rootCommand.PersistentFlags().StringVar(&proxyStr, "proxy", "", "Proxy URL (e.g. http://localhost:8080)")
}

func Execute() {
	config, err := cli.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if proxyStr != "" {
		config.Proxy = proxyStr
	}

	var transport http.RoundTripper = http.DefaultClient.Transport

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			log.Println(err)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	tokenStore, err := cli.NewFileTokenStore()
	if err != nil {
		log.Fatalln(err)
	}

	if strings.HasSuffix(config.BaseURL, "/") {
		config.BaseURL = config.BaseURL[:len(config.BaseURL)-1]
	}

	client = &cli.Client{
		Configuration: config,
		Client: &http.Client{
			Transport: transport,
		},
		TokenStore: tokenStore,
	}

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
