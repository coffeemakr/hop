package cmd

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var generateConfigCommand = &cobra.Command{
	Use:  "generate-config",
	RunE: runGenerateConfig,
}
var generateOutputPath string

type ListenConfig struct {
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
}

func (c *ListenConfig) GetServerAddress() string {
	addr := c.Host + ":" + strconv.FormatInt(int64(c.Port), 10)
	return addr
}

type DatabaseConfig struct {
	URL string `json:"url"`
}

type AuthenticationConfig struct {
	Key string `json:"key" yaml:"key,omitempty"`
}

type Configuration struct {
	Listen   *ListenConfig         `json:"listen,omitempty" yaml:",omitempty"`
	Database *DatabaseConfig       `json:"database,omitempty" yaml:",omitempty"`
	Auth     *AuthenticationConfig `json:"auth,omitempty" yaml:",omitempty"`
}

func init() {
	generateConfigCommand.PersistentFlags().StringVarP(&generateOutputPath, "output", "o", "", "The output (use - for stdout)")

}

func runGenerateConfig(*cobra.Command, []string) error {

	jwksPath, err := filepath.Abs("private-server-jwk.json")
	if err != nil {
		return err
	}
	config := Configuration{
		Listen: &ListenConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
		Database: &DatabaseConfig{
			URL: "mongodb://username:password@localhost:27017",
		},
		Auth: &AuthenticationConfig{
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
