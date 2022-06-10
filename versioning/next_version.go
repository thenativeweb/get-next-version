package versioning

import (
	"errors"

	"github.com/Masterminds/semver"
	"github.com/thenativeweb/getnextversion/conventionalcommits"
)

func CalculateNextVersion(
	currentVersion *semver.Version,
	conventionalCommitTypes []conventionalcommits.Type,
) (semver.Version, error) {
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
		return *currentVersion, nil
	case conventionalcommits.Fix:
		return currentVersion.IncPatch(), nil
	case conventionalcommits.Feature:
		return currentVersion.IncMinor(), nil
	case conventionalcommits.BreakingChange:
		return currentVersion.IncMajor(), nil
	}

	return *currentVersion, errors.New("invalid conventional commit type")
}
