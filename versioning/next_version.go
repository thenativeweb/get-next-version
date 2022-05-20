package versioning

import (
	"github.com/Masterminds/semver"
)

type SemanticVersioningTag int

const (
	Chore SemanticVersioningTag = iota
	Fix
	Feature
	BreakingChange
)

func CalculateNextVersion(currentVersion *semver.Version, semanticVersioningTags []SemanticVersioningTag) semver.Version {
	currentlyDetectedChange := Chore
	for _, tag := range semanticVersioningTags {
		if tag > currentlyDetectedChange {
			currentlyDetectedChange = tag
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

	panic("Invalid semantic versioning tag was provided. 💔")
}
