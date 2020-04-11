package main

import (
	"log"

	"github.com/coffeemakr/ruck/server/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
