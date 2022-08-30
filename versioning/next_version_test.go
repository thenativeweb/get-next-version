package versioning_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/versioning"
)

func TestCalculateNextVersion(t *testing.T) {
	tests := []struct {
		currentVersion         *semver.Version
		conventionalCommitType []git.ConventionalCommit
		expectedNewVersion     *semver.Version
		expectedHasNewVersion  bool
	}{
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{},
			expectedNewVersion:     semver.MustParse("1.0.0"),
			expectedHasNewVersion:  false,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.Chore, Message: "chore: Do something"}},
			expectedNewVersion:     semver.MustParse("1.0.0"),
			expectedHasNewVersion:  false,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.Fix, Message: "fix: Do something"}},
			expectedNewVersion:     semver.MustParse("1.0.1"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.Feature, Message: "feat: Do something"}},
			expectedNewVersion:     semver.MustParse("1.1.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.BreakingChange, Message: "feat!: Do something"}},
			expectedNewVersion:     semver.MustParse("2.0.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.1"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.Feature, Message: "feat: Do something"}},
			expectedNewVersion:     semver.MustParse("1.1.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.1.1"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.BreakingChange, Message: "feat!: Do something"}},
			expectedNewVersion:     semver.MustParse("2.0.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion: semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{
				{Type: conventionalcommits.Chore, Message: "chore: Do something"}, {Type: conventionalcommits.Fix, Message: "fix: Do something"},
				{Type: conventionalcommits.Feature, Message: "feat: Do something"}, {Type: conventionalcommits.BreakingChange, Message: "feat!: Do something"}},
			expectedNewVersion:    semver.MustParse("2.0.0"),
			expectedHasNewVersion: true,
		},
		{
			currentVersion: semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{
				{Type: conventionalcommits.BreakingChange, Message: "feat!: Do something"}, {Type: conventionalcommits.Feature, Message: "feat: Do something"},
				{Type: conventionalcommits.Fix, Message: "fix: Do something"}, {Type: conventionalcommits.Chore, Message: "chore: Do something"}},
			expectedNewVersion:    semver.MustParse("2.0.0"),
			expectedHasNewVersion: true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.Chore, Message: "chore: Do something"}, {Type: conventionalcommits.Chore, Message: "chore: Do something"}},
			expectedNewVersion:     semver.MustParse("1.0.0"),
			expectedHasNewVersion:  false,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.Fix, Message: "fix: Do something"}, {Type: conventionalcommits.Fix, Message: "fix: Do something"}},
			expectedNewVersion:     semver.MustParse("1.0.1"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.Feature, Message: "feat: Do something"}, {Type: conventionalcommits.Feature, Message: "feat: Do something"}},
			expectedNewVersion:     semver.MustParse("1.1.0"),
			expectedHasNewVersion:  true,
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			conventionalCommitType: []git.ConventionalCommit{{Type: conventionalcommits.BreakingChange, Message: "feat!: Do something"}, {Type: conventionalcommits.BreakingChange, Message: "feat!: Do something"}},
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
