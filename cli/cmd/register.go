package cmd

import (
	"bytes"
	"fmt"
	"github.com/coffeemakr/ruck"
	"github.com/spf13/cobra"
	"log"
)

var registerCommand  = &cobra.Command{
	Use: "register",
	Run: runRegister,
}

func readRegistration() (*ruck.RegistrationRequest, error) {
	var passwordsAreValid bool
	var password, passwordConfirmation []byte
	fmt.Print("Name                  :")
	name, err := readPlainText()
	if err != nil {
		return nil, err
	}

	fmt.Print("Mail                  :")
	email, err := readPlainText()
	if err != nil {
		return nil, err
	}

	for !passwordsAreValid {
		fmt.Print("Password              :")
		password, err = readPassword()
		fmt.Println()
		if err != nil {
			return nil, err
		}

		fmt.Print("Password Confirmation :")
		passwordConfirmation, err = readPassword()
		fmt.Println()
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(password, passwordConfirmation) {
			fmt.Println("\nPassword don't match. Please try again.")
			continue
		}
		passwordsAreValid = true
	}
	return &ruck.RegistrationRequest{
		Name:                 name,
		Email:                email,
		Password:             password,
		PasswordConfirmation: passwordConfirmation,
	}, nil
}

func runRegister(cmd *cobra.Command, args []string) {
	request, err := readRegistration()
	if err != nil {
		log.Fatalln(err)
	}
	user, err := client.Register(request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Registred user: %s\n", user)

}
