package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coffeemakr/wedo"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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
	client := http.Client{}
	jsonValue, _ := json.Marshal(task)
	_, err := client.Post("http://localhost:8080/tasks", "application/json", bytes.NewReader(jsonValue))
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {

}
