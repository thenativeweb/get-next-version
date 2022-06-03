package cli

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/thenativeweb/getnextversion/git"
	"github.com/thenativeweb/getnextversion/versioning"
)

var RootCommand = &cobra.Command{
	Use:   "get-next-version",
	Short: "Calculate the next version number for your project.",
	Long:  "Calculate the next version number for your project based on your git history.",
	Run: func(command *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(command.UsageString())
		}

		repository, err := gogit.PlainOpen(args[0])
		if err != nil {
			panic(err.Error())
		}

		result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
		if err != nil {
			panic(err.Error())
		}

		nextVersion := versioning.CalculateNextVersion(result.LatestReleaseVersion, result.ConventionalCommitTypes)

		fmt.Println(nextVersion.String())
	},
}
