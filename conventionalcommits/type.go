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
	defaultChoreTypes   = []string{"build", "chore", "ci", "docs", "style", "refactor", "perf", "test"}
	defaultFixTypes     = []string{"fix"}
	defaultFeatureTypes = []string{"feat"}
)

type TypeClassifier struct {
	choreTypes   []string
	fixTypes     []string
	featureTypes []string
	allTypes     []string
}

func NewTypeClassifier() *TypeClassifier {
	tc := &TypeClassifier{
		choreTypes:   append([]string{}, defaultChoreTypes...),
		fixTypes:     append([]string{}, defaultFixTypes...),
		featureTypes: append([]string{}, defaultFeatureTypes...),
	}
	tc.updateAllTypes()
	return tc
}

func NewTypeClassifierWithCustomPrefixes(customChoreTypes, customFixTypes, customFeatureTypes []string) *TypeClassifier {
	tc := &TypeClassifier{}
	
	if len(customChoreTypes) > 0 {
		tc.choreTypes = customChoreTypes
	} else {
		tc.choreTypes = append([]string{}, defaultChoreTypes...)
	}
	
	if len(customFixTypes) > 0 {
		tc.fixTypes = customFixTypes
	} else {
		tc.fixTypes = append([]string{}, defaultFixTypes...)
	}
	
	if len(customFeatureTypes) > 0 {
		tc.featureTypes = customFeatureTypes
	} else {
		tc.featureTypes = append([]string{}, defaultFeatureTypes...)
	}
	
	tc.updateAllTypes()
	return tc
}

func (tc *TypeClassifier) updateAllTypes() {
	tc.allTypes = nil
	for _, types := range [][]string{tc.choreTypes, tc.fixTypes, tc.featureTypes} {
		tc.allTypes = append(tc.allTypes, types...)
	}
}

func (tc *TypeClassifier) GetAllTypes() []string {
	return tc.allTypes
}

func (tc *TypeClassifier) StringToType(s string) (Type, error) {
	lowerS := strings.ToLower(s)
	
	if slices.Contains(tc.featureTypes, lowerS) {
		return Feature, nil
	}

	if slices.Contains(tc.fixTypes, lowerS) {
		return Fix, nil
	}

	if slices.Contains(tc.choreTypes, lowerS) {
		return Chore, nil
	}

	return Chore, errors.New("invalid string for conventional commit type")
}

var defaultClassifier = NewTypeClassifier()

func StringToType(s string) (Type, error) {
	return defaultClassifier.StringToType(s)
}

func SetCustomPrefixes(customChoreTypes, customFixTypes, customFeatureTypes []string) {
	defaultClassifier = NewTypeClassifierWithCustomPrefixes(customChoreTypes, customFixTypes, customFeatureTypes)
	initCommitMessage()
}

func initType() {
	defaultClassifier = NewTypeClassifier()
}
