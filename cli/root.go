package cli

import (
	"strings"
	
	"github.com/Masterminds/semver"
	gogit "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/target"
	"github.com/thenativeweb/get-next-version/util"
	"github.com/thenativeweb/get-next-version/versioning"
	"golang.org/x/exp/slices"
)

var (
	rootRepositoryFlag      string
	rootTargetFlag          string
	rootPrefixFlag          string
	rootFeaturePrefixesFlag string
	rootFixPrefixesFlag     string
	rootChorePrefixesFlag   string
)

func init() {
	RootCommand.Flags().StringVarP(&rootRepositoryFlag, "repository", "r", ".", "sets the path to the repository")
	RootCommand.Flags().StringVarP(&rootTargetFlag, "target", "t", "version", "sets the output target")
	RootCommand.Flags().StringVarP(&rootPrefixFlag, "prefix", "p", "", "sets the version prefix")
	RootCommand.Flags().StringVar(&rootFeaturePrefixesFlag, "feature-prefixes", "", "sets custom feature prefixes (comma-separated)")
	RootCommand.Flags().StringVar(&rootFixPrefixesFlag, "fix-prefixes", "", "sets custom fix prefixes (comma-separated)")
	RootCommand.Flags().StringVar(&rootChorePrefixesFlag, "chore-prefixes", "", "sets custom chore prefixes (comma-separated)")
}

var RootCommand = &cobra.Command{
	Use:   "get-next-version",
	Short: "Get the next version according for semantic versioning",
	Long:  "Get the next version according for semantic versioning.",
	Run: func(_ *cobra.Command, _ []string) {
		validTargets := []string{
			"github-action",
			"json",
			"version",
		}

		if isValid, prefixValidationError := util.IsValidVersionPrefix(rootPrefixFlag); !isValid {
			log.Fatal().Msgf("invalid version prefix %+q", prefixValidationError)
		}

		if !slices.Contains(validTargets, rootTargetFlag) {
			log.Fatal().Msg("invalid target")
		}

		classifier := createTypeClassifier()

		repository, err := gogit.PlainOpen(rootRepositoryFlag)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		var nextVersion semver.Version
		var hasNextVersion bool
		result, err := git.GetConventionalCommitTypesSinceLastRelease(repository, classifier)
		if err != nil {
			log.Fatal().Msg(err.Error())
		} else {
			nextVersion, hasNextVersion = versioning.CalculateNextVersion(result.LatestReleaseVersion, result.ConventionalCommitTypes)
		}

		err = target.WriteOutput(nextVersion, hasNextVersion, rootTargetFlag, rootPrefixFlag)
		if err != nil {
			log.Fatal().Err(err).Msg("could not write output")
		}
	},
}

func createTypeClassifier() *conventionalcommits.TypeClassifier {
	var choreTypes, fixTypes, featureTypes []string
	
	if rootChorePrefixesFlag != "" {
		choreTypes = parseCommaSeparatedPrefixes(rootChorePrefixesFlag)
	}
	
	if rootFixPrefixesFlag != "" {
		fixTypes = parseCommaSeparatedPrefixes(rootFixPrefixesFlag)
	}
	
	if rootFeaturePrefixesFlag != "" {
		featureTypes = parseCommaSeparatedPrefixes(rootFeaturePrefixesFlag)
	}
	
	return conventionalcommits.NewTypeClassifierWithCustomPrefixes(choreTypes, fixTypes, featureTypes)
}

func parseCommaSeparatedPrefixes(input string) []string {
	if input == "" {
		return nil
	}
	
	var result []string
	for _, prefix := range strings.Split(input, ",") {
		trimmed := strings.TrimSpace(prefix)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
