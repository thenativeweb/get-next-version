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
	choreTypes   = []string{"build", "chore", "ci", "docs", "style", "refactor", "perf", "test"}
	fixTypes     = []string{"fix"}
	featureTypes = []string{"feat"}
	allTypes     []string
)

func initType() {
	for _, types := range [][]string{choreTypes, fixTypes, featureTypes} {
		allTypes = append(allTypes, types...)
	}
}

func StringToType(s string) (Type, error) {
	if slices.Contains(choreTypes, strings.ToLower(s)) {
		return Chore, nil
	}

	if slices.Contains(fixTypes, strings.ToLower(s)) {
		return Fix, nil
	}

	if slices.Contains(featureTypes, strings.ToLower(s)) {
		return Feature, nil
	}

	return Chore, errors.New("invalid string for conventional commit type")
}
