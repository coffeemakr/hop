package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
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
	Aliases: []string{"ls"},
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

var groupSetDefaultCommand = &cobra.Command{
	Use: "set-default",
	Run: runSetDefaultGroup,
	Args: cobra.ExactArgs(1),
}

var groupGetDefaultCommand = &cobra.Command{
	Use: "get-default",
	Run: runGetDefaultGroup,
	Args: cobra.NoArgs,
}

func runSetDefaultGroup(cmd *cobra.Command, args []string) {
	err := setDefaultGroup(args[0])
	if err != nil {
		log.Fatalln(err)
	}
}

func runGetDefaultGroup(cmd *cobra.Command, args []string) {
	fmt.Printf("Default Group ID: %s\n", getDefaultGroup())
}

func getDefaultGroup() string {
	return viper.GetString("default-group")
}

func viperWriteOrCreateConfig() (err error){
	const defaultConfigurationFile = "$HOME/.config/wedo-config.yml"
	path := os.ExpandEnv(defaultConfigurationFile)

	fp, err := os.Create(path)
	if err != nil {
		return
	}
	err = fp.Close()
	if err != nil {
		return
	}
	return viper.WriteConfig()
}

func setDefaultGroup(groupId string) error {
	viper.Set("default-group", groupId)
	return viperWriteOrCreateConfig()
}

func runAddGroup(cmd *cobra.Command, args []string) {
	groupName := args[0]
	err := client.CreateGroup(groupName)
	if err != nil {
		log.Fatalln(err)
	}
	defaultGroup := getDefaultGroup()
	if defaultGroup == "" {
		err := setDefaultGroup(defaultGroup)
		err = fmt.Errorf("failed to set default group: %s", err)
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
		setDefaultGroup("")
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
	groupCommand.AddCommand(groupAddCommand, groupListCommand, groupPruneCommand, groupJoinCommand, groupSetDefaultCommand, groupGetDefaultCommand)
}
