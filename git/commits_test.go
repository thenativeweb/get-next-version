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
				{message: "chore: Do something", tag: "1.0.0"},
				{message: "chore: Do something else", tag: ""},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{conventionalcommits.Chore},
		},
	}

	for _, test := range tests {
		repository := testutil.SetUpInMemoryRepository()

		for _, commit := range test.commitHistory {
			worktree, _ := repository.Worktree()
			worktree.Commit(commit.message, testutil.CommitOptions)

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
		assert.Equal(t, test.expectedConventionalCommitTypes, actual.ConventionalCommitTypes)
	}

}
