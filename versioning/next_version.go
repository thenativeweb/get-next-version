package versioning

import (
	"github.com/Masterminds/semver"
	"github.com/thenativeweb/getnextversion/conventionalcommits"
)

func CalculateNextVersion(
	currentVersion *semver.Version,
	conventionalCommitTypes []conventionalcommits.ConventionalCommitType,
) semver.Version {
	currentlyDetectedChange := conventionalcommits.Chore
	for _, commitType := range conventionalCommitTypes {
		if commitType > currentlyDetectedChange {
			currentlyDetectedChange = commitType
		}
		if currentlyDetectedChange == conventionalcommits.BreakingChange {
			break
		}
	}

	switch currentlyDetectedChange {
	case conventionalcommits.Chore:
		return *currentVersion
	case conventionalcommits.Fix:
		return currentVersion.IncPatch()
	case conventionalcommits.Feature:
		return currentVersion.IncMinor()
	case conventionalcommits.BreakingChange:
		return currentVersion.IncMajor()
	}

	panic("Invalid conventional commit type")
}
