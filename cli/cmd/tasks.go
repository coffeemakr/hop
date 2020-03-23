package cmd

import (
	"fmt"
	"github.com/coffeemakr/wedo"
	"github.com/spf13/cobra"
	"log"
	"os"
	"text/template"
)

var (
	taskCommand = &cobra.Command{
		Use: "task",
	}
	taskListCommand = &cobra.Command{
		Use: "list",
		Run: runTaskList,
		Aliases: []string{"ls"},
	}
	taskGetCommand = &cobra.Command{
		Use:  "get",
		Run:  runTaskGet,
		Args: cobra.ExactArgs(1),
	}
	taskAddCommand = &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run:   runAddTask,
		Args:  cobra.ExactArgs(1),
	}
)

func runTaskGet(cmd *cobra.Command, args []string) {
	taskId := args[0]
	var task *wedo.Task
	task, err := client.GetTaskDetails(taskId)
	if err != nil {
		log.Fatalln(err)
	}
	t, err := template.New("taskTemplate").Parse("ID    {{.ID}}\n" +
		"Name  {{.Name}}\n" +
		"Group {{.Group}}\n")
	if err != nil {
		log.Fatalln(err)
	}
	err = t.Execute(os.Stdout, task)
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	taskCommand.AddCommand(taskAddCommand, taskListCommand, taskGetCommand)
}

func runTaskList(cmd *cobra.Command, args []string) {
	tasks, err := client.GetTaskList()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, task := range tasks {
		fmt.Printf("Task %s: %-40s %-20s %s\n", task.ID, task.Name, task.AssigneeName, task.DueDate.Format("01.02.2006"))
	}
}

func runAddTask(cmd *cobra.Command, args []string) {
	var name = args[0]
	var task = wedo.NewTask(name, wedo.Weeks, 1)
	var defaultGroupId = getDefaultGroup()
	if defaultGroupId == "" {
		log.Fatalln("no default group set")
	}
	err := client.CreateTaskForGroup(task, defaultGroupId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Task created: %s\n", task)

}
