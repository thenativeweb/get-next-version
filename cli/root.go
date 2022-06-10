package cli

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thenativeweb/getnextversion/git"
	"github.com/thenativeweb/getnextversion/versioning"
)

var RootCommand = &cobra.Command{
	Use:   "get-next-version",
	Short: "Get the next semantic version for your project",
	Long:  "Get the next semantic version for your project based on your git history.",
	Run: func(command *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(command.UsageString())
			return
		}

		repository, err := gogit.PlainOpen(args[0])
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
		if err != nil {
			if err == git.NoCommitsFoundError {
				fmt.Println("0.0.1")
				return
			}
			log.Fatal().Msg(err.Error())
		}

		nextVersion, err := versioning.CalculateNextVersion(result.LatestReleaseVersion, result.ConventionalCommitTypes)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		fmt.Println(nextVersion.String())
	},
}
