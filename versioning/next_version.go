package versioning

import (
	"github.com/Masterminds/semver"
)

type ConventionalCommitType int

const (
	Chore ConventionalCommitType = iota
	Fix
	Feature
	BreakingChange
)

func CalculateNextVersion(currentVersion *semver.Version, conventionalCommitTypes []ConventionalCommitType) semver.Version {
	currentlyDetectedChange := Chore
	for _, commitType := range conventionalCommitTypes {
		if commitType > currentlyDetectedChange {
			currentlyDetectedChange = commitType
		}
		if currentlyDetectedChange == BreakingChange {
			break
		}
	}

	switch currentlyDetectedChange {
	case Chore:
		return *currentVersion
	case Fix:
		return currentVersion.IncPatch()
	case Feature:
		return currentVersion.IncMinor()
	case BreakingChange:
		return currentVersion.IncMajor()
	}

	panic("Invalid conventional commit type")
}
