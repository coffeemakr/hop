package cmd

import "github.com/spf13/cobra"

var groupComand = &cobra.Command{
	Use: "group",
}

var groupAddCommand = &cobra.Command{
	Use: "add",
}

func init() {
	groupComand.AddCommand()
}
