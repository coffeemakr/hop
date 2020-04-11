package main

import (
	"log"

	"github.com/coffeemakr/amtli/server/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
