package cli

import (
	"fmt"

	"github.com/Masterminds/semver"
	gogit "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thenativeweb/get-next-version/cliutil"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/versioning"
	"golang.org/x/exp/slices"
)

var rootRepositoryFlag string
var rootFormatFlag string

func init() {
	RootCommand.Flags().StringVarP(&rootRepositoryFlag, "repository", "r", ".", "sets the path to the repository")
	RootCommand.Flags().StringVarP(&rootFormatFlag, "format", "f", "version", "sets the output format")
}

var RootCommand = &cobra.Command{
	Use:   "get-next-version",
	Short: "Get the next version according for semantic versioning",
	Long:  "Get the next version according for semantic versioning.",
	Run: func(_ *cobra.Command, _ []string) {
		validFormats := []string{
			"github-action",
			"json",
			"version",
		}

		if !slices.Contains(validFormats, rootFormatFlag) {
			log.Fatal().Msg("invalid format")
		}

		repository, err := gogit.PlainOpen(rootRepositoryFlag)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		var nextVersion semver.Version
		var hasNextVersion bool
		result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
		if err != nil {
			log.Fatal().Msg(err.Error())
		} else {
			nextVersion, hasNextVersion = versioning.CalculateNextVersion(result.LatestReleaseVersion, result.ConventionalCommitTypes)
		}

		lines := cliutil.Format(nextVersion, hasNextVersion, rootFormatFlag)
		for _, line := range lines {
			fmt.Println(line)
		}
	},
}
