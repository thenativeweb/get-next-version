package main

import (
	"log"

	"github.com/thenativeweb/getnextversion/cli"
)

func main() {
	err := cli.RootCommand.Execute()
	if err != nil {
		log.Fatal("failed to execute root command")
	}
}
