package conventionalcommits

import (
	"errors"
	"strings"
)

type Type int

const (
	Chore Type = iota
	Fix
	Feature
	BreakingChange
)

var choreTypes = []string{"build", "chore", "ci", "docs", "style", "refector", "perf", "test"}
var fixTypes = []string{"fix"}
var featureTypes = []string{"feat"}
var allTypes []string

func initType() {
	for _, types := range [][]string{choreTypes, fixTypes, featureTypes} {
		allTypes = append(allTypes, types...)
	}
}

func StringToType(s string) (Type, error) {
	for _, choreType := range choreTypes {
		if strings.ToLower(s) == choreType {
			return Chore, nil
		}
	}

	for _, fixType := range fixTypes {
		if strings.ToLower(s) == fixType {
			return Fix, nil
		}
	}

	for _, featureType := range featureTypes {
		if strings.ToLower(s) == featureType {
			return Feature, nil
		}
	}

	return Chore, errors.New("invalid string for conventional commit type")
}
