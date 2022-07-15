package cli

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/versioning"
)

var rootRepositoryFlag string

func init() {
	RootCommand.Flags().StringVarP(&rootRepositoryFlag, "repository", "r", ".", "path of the repository to operate on")
}

var RootCommand = &cobra.Command{
	Use:   "get-next-version",
	Short: "Get the next semantic version for your project",
	Long:  "Get the next semantic version for your project based on your git history.",
	Run: func(command *cobra.Command, _ []string) {
		repository, err := gogit.PlainOpen(rootRepositoryFlag)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
		if err != nil {
			if err == git.ErrNoCommitsFound {
				fmt.Println("0.0.1")
				return
			}
			log.Fatal().Msg(err.Error())
		}

		nextVersion := versioning.CalculateNextVersion(result.LatestReleaseVersion, result.ConventionalCommitTypes)
		fmt.Println(nextVersion.String())
	},
}
