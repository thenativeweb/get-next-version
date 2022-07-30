package git_test

import (
	"testing"

	"github.com/Masterminds/semver"
	gogit "github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/testutil"
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
		annotateTags                    bool
	}{
		{
			commitHistory:                   []commit{},
			doExpectError:                   true,
			expectedLastVersion:             nil,
			expectedConventionalCommitTypes: []conventionalcommits.Type{},
			annotateTags:                    false,
		},
		{
			commitHistory: []commit{
				{message: "chore: Do something", tag: ""},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("0.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{conventionalcommits.Chore},
			annotateTags:                    false,
		},
		{
			commitHistory: []commit{
				{message: "Last release", tag: "1.0.0"},
				{message: "Do something", tag: ""},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{conventionalcommits.Chore},
			annotateTags:                    false,
		},
		{
			commitHistory: []commit{
				{message: "chore: Do something", tag: "1.0.0"},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{},
			annotateTags:                    false,
		},
		{
			commitHistory: []commit{
				{message: "chore: irrelevant", tag: "0.0.1"},
				{message: "feat: because it is", tag: ""},
				{message: "feat(scope)!: before the last tag", tag: "0.0.2"},
				{message: "chore: Do something", tag: "1.0.0"},
				{message: "chore: Do something else", tag: ""},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{conventionalcommits.Chore},
			annotateTags:                    false,
		},
		{
			commitHistory: []commit{
				{message: "chore: irrelevant", tag: "v0.0.1"},
				{message: "feat: because it is", tag: ""},
				{message: "feat(scope)!: before the last tag", tag: "0.0.2"},
				{message: "chore: Do something", tag: "v1.0.0"},
				{message: "chore: Do something else", tag: ""},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{conventionalcommits.Chore},
			annotateTags:                    false,
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
			annotateTags: false,
		},
		{
			commitHistory: []commit{
				{message: "Last release", tag: "1.0.0"},
				{message: "fix: Do something", tag: ""},
			},
			doExpectError:                   false,
			expectedLastVersion:             semver.MustParse("1.0.0"),
			expectedConventionalCommitTypes: []conventionalcommits.Type{conventionalcommits.Fix},
			annotateTags:                    true,
		},
	}

	for _, test := range tests {
		repository := testutil.SetUpInMemoryRepository()
		worktree, _ := repository.Worktree()

		for _, commit := range test.commitHistory {
			commitOptions := testutil.CreateCommitOptions()
			worktree.Commit(commit.message, commitOptions)

			if commit.tag == "" {
				continue
			}

			head, _ := repository.Head()
			var createTagOpts *gogit.CreateTagOptions
			if test.annotateTags {
				createTagOpts = &gogit.CreateTagOptions{
					Message: "some message",
					Tagger:  commitOptions.Author,
				}
			}
			repository.CreateTag(commit.tag, head.Hash(), createTagOpts)
		}

		actual, err := git.GetConventionalCommitTypesSinceLastRelease(repository)

		if test.doExpectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)

		// We need to compare using semver's own Equal function, because under
		// the hood it makes a difference between versions with a leading v and
		// versions without a leading v. However, when printing them, they are
		// actually both shown without the v.
		assert.True(t, test.expectedLastVersion.Equal(actual.LatestReleaseVersion))
		assert.Equal(t, test.expectedLastVersion.String(), actual.LatestReleaseVersion.String())
		assert.ElementsMatch(t, test.expectedConventionalCommitTypes, actual.ConventionalCommitTypes)
	}
}
