package cmd

import "github.com/spf13/cobra"

const (
	jwkName = "private-server-jwk.json"
)

var rootCmd = &cobra.Command{
	Use: "hop-server",
}

func init() {
	rootCmd.AddCommand(generateKeysCommand, generateConfigCommand, serverCommand)
}

func Execute() error {
	return rootCmd.Execute()
}
