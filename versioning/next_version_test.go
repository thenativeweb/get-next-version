package versioning_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thenativeweb/get-next-version/versioning"
)

func TestCalculateNextVersion(t *testing.T) {
	tests := []struct {
		currentVersion         *semver.Version
		conventionalCommitType []conventionalcommits.Type
		expectedNewVersion     *semver.Version
		expectedHasNewVersion  bool
	}{
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{},
			expectedNewVersion:     semver.MustParse("1.0.0"),
			expectedHasNewVersion:  false,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
			expectedHasNewVersion:  false,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.1"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.1.1"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Chore, conventionalcommits.Fix, conventionalcommits.Feature, conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.BreakingChange, conventionalcommits.Feature, conventionalcommits.Fix, conventionalcommits.Chore},
			expectedNewVersion:     semver.MustParse("2.0.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Chore, conventionalcommits.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
			expectedHasNewVersion:  false,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Fix, conventionalcommits.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.Feature, conventionalcommits.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []conventionalcommits.Type{conventionalcommits.BreakingChange, conventionalcommits.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
			expectedHasNewVersion:  true,
		},
	}

	for _, test := range tests {
		actualNewVersion, hasNewVersion := versioning.CalculateNextVersion(test.currentVersion, test.conventionalCommitType)
		assert.Equal(t, test.expectedNewVersion, &actualNewVersion)
		assert.Equal(t, test.expectedHasNewVersion, hasNewVersion)
	}
}
