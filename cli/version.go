package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thenativeweb/get-next-version/version"
)

func init() {
	RootCommand.AddCommand(VersionCommand)
}

var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the get-next-version version",
	Long:  "Prints the get-next-version version.",
	Run: func(command *cobra.Command, args []string) {
		fmt.Println("get-next-version " + version.Version)
		fmt.Println("Revision " + version.GitVersion)
	},
}
