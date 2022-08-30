package versioning

import (
	"github.com/Masterminds/semver"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thenativeweb/get-next-version/git"
)

func CalculateNextVersion(
	currentVersion *semver.Version,
	conventionalCommits []git.ConventionalCommit,
) (semver.Version, bool) {
	currentlyDetectedChange := conventionalcommits.Chore
	for _, commit := range conventionalCommits {
		if commit.Type > currentlyDetectedChange {
			currentlyDetectedChange = commit.Type
		}
		if currentlyDetectedChange == conventionalcommits.BreakingChange {
			break
		}
	}

	switch currentlyDetectedChange {
	case conventionalcommits.Chore:
		return *currentVersion, false
	case conventionalcommits.Fix:
		return currentVersion.IncPatch(), true
	case conventionalcommits.Feature:
		return currentVersion.IncMinor(), true
	case conventionalcommits.BreakingChange:
		return currentVersion.IncMajor(), true
	}

	panic("invalid conventional commit type")
}
