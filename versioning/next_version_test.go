package versioning_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/getnextversion/conventionalcommits"
	"github.com/thenativeweb/getnextversion/versioning"
)

func TestCalculateNextVersion(t *testing.T) {
	tests := []struct {
		currentVersion         *semver.Version
		conventionalCommitType []conventionalcommits.ConventionalCommitType
		expectedNewVersion     *semver.Version
	}{
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.1"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.1.1"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Chore, conventionalcommits.Fix, conventionalcommits.Feature, conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.BreakingChange, conventionalcommits.Feature, conventionalcommits.Fix, conventionalcommits.Chore},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Chore, conventionalcommits.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Fix, conventionalcommits.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.Feature, conventionalcommits.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.ConventionalCommitType{conventionalcommits.BreakingChange, conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
	}

	for _, test := range tests {
		actualNewVersion := versioning.CalculateNextVersion(test.currentVersion, test.conventionalCommitType)
		assert.EqualValues(t, test.expectedNewVersion, &actualNewVersion)
	}
}
