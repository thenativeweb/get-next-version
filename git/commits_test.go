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
		commitHistory               []commit
		doExpectError               bool
		expectedLastVersion         *semver.Version
		expectedConventionalCommits []git.ConventionalCommit
		annotateTags                bool
	}{
		{
			commitHistory:               []commit{},
			doExpectError:               true,
			expectedLastVersion:         nil,
			expectedConventionalCommits: []git.ConventionalCommit{},
			annotateTags:                false,
		},
		{
			commitHistory: []commit{
				{message: "chore: Do something", tag: ""},
			},
			doExpectError:       false,
			expectedLastVersion: semver.MustParse("0.0.0"),
			expectedConventionalCommits: []git.ConventionalCommit{
				{conventionalcommits.Chore, "chore: Do something"},
			},
			annotateTags: false,
		},
		{
			commitHistory: []commit{
				{message: "Last release", tag: "1.0.0"},
				{message: "Do something", tag: ""},
			},
			doExpectError:       false,
			expectedLastVersion: semver.MustParse("1.0.0"),
			expectedConventionalCommits: []git.ConventionalCommit{
				{conventionalcommits.Chore, "Do something"},
			},
			annotateTags: false,
		},
		{
			commitHistory: []commit{
				{message: "chore: Do something", tag: "1.0.0"},
			},
			doExpectError:               false,
			expectedLastVersion:         semver.MustParse("1.0.0"),
			expectedConventionalCommits: []git.ConventionalCommit{},
			annotateTags:                false,
		},
		{
			commitHistory: []commit{
				{message: "chore: irrelevant", tag: "0.0.1"},
				{message: "feat: because it is", tag: ""},
				{message: "feat(scope)!: before the last tag", tag: "0.0.2"},
				{message: "chore: Do something", tag: "1.0.0"},
				{message: "chore: Do something else", tag: ""},
			},
			doExpectError:       false,
			expectedLastVersion: semver.MustParse("1.0.0"),
			expectedConventionalCommits: []git.ConventionalCommit{
				{conventionalcommits.Chore, "chore: Do something else"},
			},
			annotateTags: false,
		},
		{
			commitHistory: []commit{
				{message: "chore: irrelevant", tag: "v0.0.1"},
				{message: "feat: because it is", tag: ""},
				{message: "feat(scope)!: before the last tag", tag: "0.0.2"},
				{message: "chore: Do something", tag: "v1.0.0"},
				{message: "chore: Do something else", tag: ""},
			},
			doExpectError:       false,
			expectedLastVersion: semver.MustParse("1.0.0"),
			expectedConventionalCommits: []git.ConventionalCommit{
				{conventionalcommits.Chore, "chore: Do something else"},
			},
			annotateTags: false,
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
			expectedConventionalCommits: []git.ConventionalCommit{
				{conventionalcommits.Chore, "chore: non breaking"},
				{conventionalcommits.Fix, "fix: non breaking"},
				{conventionalcommits.Feature, "feat: non breaking"},
				{conventionalcommits.BreakingChange, "chore!: breaking"},
				{conventionalcommits.BreakingChange, "fix(with scope)!: breaking"},
				{conventionalcommits.BreakingChange, "feat: breaking\n\nBREAKING-CHANGE: with footer"},
			},
			annotateTags: false,
		},
		{
			commitHistory: []commit{
				{message: "Last release", tag: "1.0.0"},
				{message: "fix: Do something", tag: ""},
			},
			doExpectError:       false,
			expectedLastVersion: semver.MustParse("1.0.0"),
			expectedConventionalCommits: []git.ConventionalCommit{
				{conventionalcommits.Fix, "fix: Do something"},
			},
			annotateTags: true,
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

		actual, err := git.GetConventionalCommitsSinceLastRelease(repository)

		if test.doExpectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)

		// The test in the next line is not optimal. We rely on the Equal
		// function of the SemVer module here, which considers v1.0.0 and
		// 1.0.0 to be the same. In contrast to this, assert.Equal fails
		// when comparing these two versions, due to the leading v.
		assert.True(t, test.expectedLastVersion.Equal(actual.LatestReleaseVersion))
		assert.ElementsMatch(t, test.expectedConventionalCommits, actual.ConventionalCommits)
	}
}
