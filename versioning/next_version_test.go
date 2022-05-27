package versioning_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/getnextversion/versioning"
)

func TestCalculateNextVersion(t *testing.T) {
	tests := []struct {
		currentVersion         *semver.Version
		conventionalCommitType []versioning.ConventionalCommitType
		expectedNewVersion     *semver.Version
	}{
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.1"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.1.1"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Chore, versioning.Fix, versioning.Feature, versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.BreakingChange, versioning.Feature, versioning.Fix, versioning.Chore},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Chore, versioning.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Fix, versioning.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.Feature, versioning.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []versioning.ConventionalCommitType{versioning.BreakingChange, versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
	}

	for _, test := range tests {

		actualNewVersion := versioning.CalculateNextVersion(test.currentVersion, test.conventionalCommitType)
		assert.Equal(t, test.expectedNewVersion, &actualNewVersion)
	}

}
