package main

import (
	"os"

	"github.com/thenativeweb/get-next-version/cli"
)

func main() {
	// Cobra already prints the error to stderr, so there is nothing left to
	// report here.
	if err := cli.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
