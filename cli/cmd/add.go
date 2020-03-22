package cmd

import (
	"fmt"
	"github.com/coffeemakr/wedo"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use: "add",
	Short: "Add a new task",
	Run: runAddTask,
	Args: cobra.ExactArgs(1),
}



func runAddTask(cmd *cobra.Command, args []string) {
	var name = args[0]
	var task = wedo.NewTask(name, wedo.Weeks, 1)
	fmt.Printf("Creating task: %s\n", task)
	client.CreateTask(task)
}

func init() {

}
