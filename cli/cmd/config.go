package cmd

import (
	"github.com/coffeemakr/ruck/cli"
	"github.com/spf13/cobra"
)

var (
	configCommand = &cobra.Command{
		Use:  "config",
		RunE: runConfig,
	}
)

func runConfig(cmd *cobra.Command, args []string) error {
	config, err := cli.LoadConfig()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		var key = args[0]
		return cli.PrintConfig(config, key)
	} else if len(args) == 2 {
		var key = args[0]
		var value = args[1]
		err = cli.SetConfig(config, key, value)
		if err != nil {
			return err
		}
		err = cli.WriteConfig(config)
		if err != nil {
			return err
		}
	}
	return nil
}
