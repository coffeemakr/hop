package cmd

import (
	"errors"
	"fmt"
	"github.com/coffeemakr/amtli"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"syscall"
)

var loginCommand = &cobra.Command{
	Use: "login",
	Run: runLogin,
}

func readPassword() (password []byte, err error) {
	password, err = terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return
	}
	if len(password) > 256 {
		err = errors.New("password is too long (more than 256 characters)")
		password = nil
	}
	return
}

func readPlainText() (text string, err error) {
	_, err = fmt.Scanln(&text)
	return
}

func readCredentials() (*amtli.Credentials, error) {
	var credentials amtli.Credentials
	var err error
	fmt.Print("Name    : ")
	credentials.Name, err = readPlainText()
	if err != nil {
		return nil, err
	}
	fmt.Print("Password: ")
	credentials.Password, err = readPassword()
	fmt.Println()
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

func runLogin(cmd *cobra.Command, args []string) {
	creds, err := readCredentials()
	if err != nil {
		log.Fatalln(err)
	}
	if err := client.Login(creds); err != nil {
		log.Fatalln(err)
	}
}
