package versioning

import (
	"github.com/Masterminds/semver"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
)

func CalculateNextVersion(
	currentVersion *semver.Version,
	conventionalCommitTypes []conventionalcommits.Type,
) (semver.Version, bool) {
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
