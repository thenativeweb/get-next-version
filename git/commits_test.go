package git_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/getnextversion/conventionalcommits"
	"github.com/thenativeweb/getnextversion/git"
	"github.com/thenativeweb/getnextversion/testutil"
)

type commit struct {
	message string
	tag     string
}

func TestGetConventionalCommitTypesSinceLatestRelease(t *testing.T) {
	tests := []struct {
		commitHistory                   []commit
		doExpectError                   bool
		expectedLastVersion             *semver.Version
		expectedConventionalCommitTypes []conventionalcommits.Type
	}{
		{
			commitHistory:                   []commit{},
			doExpectError:                   true,
			expectedLastVersion:             nil,
			expectedConventionalCommitTypes: []conventionalcommits.Type{},
		},
		{
			commitHistory: []commit{
				{message: "chore: Do something", tag: ""},
			},
			doExpectError:                   true,
			expectedLastVersion:             nil,
			expectedConventionalCommitTypes: []conventionalcommits.Type{},
		},
		{
			commitHistory: []commit{
				{message: "chore: Do something", tag: "1.0.0"},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{},
		},
		{
			commitHistory: []commit{
				{message: "chore: irelevant", tag: "0.0.1"},
				{message: "feat: because it is", tag: ""},
				{message: "feat(scope)!: before the last tag", tag: "0.0.2"},
				{message: "chore: Do something", tag: "1.0.0"},
				{message: "chore: Do something else", tag: ""},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{conventionalcommits.Chore},
		},
		{
			commitHistory: []commit{
				{message: "chore: Do something", tag: "1.0.0"},
				{message: "chore: non breaking", tag: ""},
				{message: "fix: non breaking", tag: ""},
				{message: "feat: non breaking", tag: ""},
				{message: "chore!: breaking", tag: ""},
				{message: "fix(with scope)!: breaking", tag: ""},
				{message: "feat: breaking\n\nBREAKING-CHANGE: with footer", tag: ""},
			},
			doExpectError:       false,
			expectedLastVersion: semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{
				conventionalcommits.Chore,
				conventionalcommits.Fix,
				conventionalcommits.Feature,
				conventionalcommits.BreakingChange,
				conventionalcommits.BreakingChange,
				conventionalcommits.BreakingChange,
			},
		},
	}

	for _, test := range tests {
		repository := testutil.SetUpInMemoryRepository()
		worktree, _ := repository.Worktree()

		for _, commit := range test.commitHistory {
			worktree.Commit(commit.message, testutil.CreateCommitOptions())

			if commit.tag == "" {
				continue
			}

			head, _ := repository.Head()
			repository.CreateTag(commit.tag, head.Hash(), nil)
		}

		actual, err := git.GetConventionalCommitTypesSinceLastRelease(repository)

		if test.doExpectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, test.expectedLastVersion, actual.LatestReleaseVersion)
		assert.ElementsMatch(t, test.expectedConventionalCommitTypes, actual.ConventionalCommitTypes)
	}

}
