package versioning

import (
	"github.com/Masterminds/semver"
)

type semanticVersioningTag int

const (
	Chore semanticVersioningTag = iota
	Fix
	Feature
	BreakingChange
)

func CalculateNextVersion(currentVersion semver.Version, semanticVersioningTags []semanticVersioningTag) semver.Version {
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
		return currentVersion
	case Fix:
		return currentVersion.IncPatch()
	case Feature:
		return currentVersion.IncMinor()
	case BreakingChange:
		return currentVersion.IncMajor()
	}

	panic("Invalid semantic versioning tag was provided. 💔")
}
