package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var groupCommand = &cobra.Command{
	Use: "group",
}

var groupAddCommand = &cobra.Command{
	Use:  "add",
	Run:  runAddGroup,
	Args: cobra.ExactArgs(1),
}

var groupListCommand = &cobra.Command{
	Use: "list",
	Run: runListGroup,
}

var groupPruneCommand = &cobra.Command{
	Use: "prune",
	Run: runPruneGroup,
}

var groupJoinCommand = &cobra.Command{
	Use: "join",
	Run: runJoinGroup,
	Args: cobra.ExactArgs(1),
}

func runAddGroup(cmd *cobra.Command, args []string) {
	groupName := args[0]
	err := client.CreateGroup(groupName)
	if err != nil {
		log.Fatalln(err)
	}
}

func runListGroup(cmd *cobra.Command, args []string) {
	groups, err := client.ListGroup()
	if err != nil {
		log.Fatalln(err)
	}
	if len(groups) == 0 {
		fmt.Println("No groups.")
	}
	for _, group := range groups {
		fmt.Printf("Group %s: %s\n", group.ID, group.Name)
	}
}

func runPruneGroup(cmd *cobra.Command, args []string) {
	groups, err := client.ListGroup()
	if err != nil {
		log.Fatalln(err)
	}
	for _, group := range groups {
		err := client.DeleteGroupByID(group.ID)
		if err != nil {
			log.Println(err)
		}
	}
}


func runJoinGroup(cmd *cobra.Command, args []string) {
	groupId := args[0]
	err := client.JoinGroup(groupId)
	if err != nil {
		log.Fatalln(err)
	}
}


func init() {
	groupCommand.AddCommand(groupAddCommand, groupListCommand, groupPruneCommand, groupJoinCommand)
}
