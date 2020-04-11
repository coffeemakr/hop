package cmd

import (
	"os"
	"path/filepath"

	"github.com/coffeemakr/ruck/server"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var generateConfigCommand = &cobra.Command{
	Use:  "generate-config",
	RunE: runGenerateConfig,
}
var generateOutputPath string

func init() {
	generateConfigCommand.PersistentFlags().StringVarP(&generateOutputPath, "output", "o", "", "The output (use - for stdout)")

}

func runGenerateConfig(*cobra.Command, []string) error {

	jwksPath, err := filepath.Abs("private-server-jwk.json")
	if err != nil {
		return err
	}
	config := server.Configuration{
		Listen: &server.ListenConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
		Database: &server.DatabaseConfig{
			URL: "mongodb://username:password@localhost:27017",
		},
		Auth: &server.AuthenticationConfig{
			Key: jwksPath,
		},
	}
	encoder := yaml.NewEncoder(os.Stdout)
	//encoder := json.NewEncoder(os.Stdout)
	//encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func init() {
}
