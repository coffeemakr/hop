package cmd

import "github.com/spf13/cobra"

var doneCommand = &cobra.Command{
	Use: "done",
	RunE: runFinishTask,
}

func runFinishTask(cmd *cobra.Command, args []string) error {
	return nil;
}
