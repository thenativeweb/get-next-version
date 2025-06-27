package conventionalcommits

import (
	"errors"
	"slices"
	"strings"
)

type Type int

const (
	Chore Type = iota
	Fix
	Feature
	BreakingChange
)

var (
	// Default prefix mappings
	defaultChoreTypes   = []string{"build", "chore", "ci", "docs", "style", "refactor", "perf", "test"}
	defaultFixTypes     = []string{"fix"}
	defaultFeatureTypes = []string{"feat"}
	
	// Active prefix mappings (can be overridden)
	choreTypes   = []string{"build", "chore", "ci", "docs", "style", "refactor", "perf", "test"}
	fixTypes     = []string{"fix"}
	featureTypes = []string{"feat"}
	allTypes     []string
)

func initType() {
	updateAllTypes()
}

func updateAllTypes() {
	allTypes = nil
	for _, types := range [][]string{choreTypes, fixTypes, featureTypes} {
		allTypes = append(allTypes, types...)
	}
	// Re-initialize commit message regex with updated types
	initCommitMessage()
}

// SetCustomPrefixes allows overriding the default prefix mappings
func SetCustomPrefixes(customChoreTypes, customFixTypes, customFeatureTypes []string) {
	if len(customChoreTypes) > 0 {
		choreTypes = customChoreTypes
	} else {
		choreTypes = append([]string{}, defaultChoreTypes...)
	}
	
	if len(customFixTypes) > 0 {
		fixTypes = customFixTypes
	} else {
		fixTypes = append([]string{}, defaultFixTypes...)
	}
	
	if len(customFeatureTypes) > 0 {
		featureTypes = customFeatureTypes
	} else {
		featureTypes = append([]string{}, defaultFeatureTypes...)
	}
	
	updateAllTypes()
}

// ResetToDefaults resets all prefix mappings to their default values
func ResetToDefaults() {
	choreTypes = append([]string{}, defaultChoreTypes...)
	fixTypes = append([]string{}, defaultFixTypes...)
	featureTypes = append([]string{}, defaultFeatureTypes...)
	updateAllTypes()
}

func StringToType(s string) (Type, error) {
	lowerS := strings.ToLower(s)
	
	// Check in order of precedence: Feature > Fix > Chore
	if slices.Contains(featureTypes, lowerS) {
		return Feature, nil
	}

	if slices.Contains(fixTypes, lowerS) {
		return Fix, nil
	}

	if slices.Contains(choreTypes, lowerS) {
		return Chore, nil
	}

	return Chore, errors.New("invalid string for conventional commit type")
}
