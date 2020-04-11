package cmd

import (
	"bufio"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"crypto/rand"

	"github.com/spf13/cobra"
	"github.com/square/go-jose/v3"
)

var (
	generateKeysCommand *cobra.Command = &cobra.Command{
		Use:  "generate-secrets",
		RunE: generateKeys,
	}
)

func generateSignatureRsaKey() (*jose.JSONWebKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	priv := jose.JSONWebKey{
		Key:       key,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}

	thumb, err := priv.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	priv.KeyID = base64.RawURLEncoding.EncodeToString(thumb)
	return &priv, nil
}

func generateKeys(*cobra.Command, []string) error {
	priv, err := generateSignatureRsaKey()
	if err != nil {
		return err
	}
	privJSON, err := priv.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = os.Stat(jwkName)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Failed to check if file exists: %s\n", err)
	}
	if err == nil {
		answer, err := askYesOrNo(fmt.Sprintf("File %s already exists.\nOverwrite?", jwkName))
		if err != nil {
			return err
		}
		if !answer {
			fmt.Printf("User abort!\n")
			os.Exit(1)
		}
	}
	err = ioutil.WriteFile(jwkName, privJSON, 0600)
	if err == nil {
		fmt.Printf("Key ID: %s\nKey successfully generated.\n", priv.KeyID)
		os.Exit(0)
	}
	return err
}

func askYesOrNo(question string) (bool, error) {
	var valid, answer bool
	for !valid {
		print(question)
		print(" [y/n] ")
		reader := bufio.NewReader(os.Stdin)
		answerString, _ := reader.ReadString('\n')
		answerString = answerString[:len(answerString)-1]
		switch answerString {
		case "y", "Y":
			answer = true
			valid = true
		case "n", "N":
			answer = false
			valid = true
		default:
			fmt.Printf("\nInvalid answer. Type 'y' or 'n'.\n")
		}
	}
	return answer, nil
}
