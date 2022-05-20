package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "get-next-version",
	Short: "Calculate the next version number for your project.",
	Long:  "Calculate the next version number for your project based on your git history.",
	Run: func(command *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}
