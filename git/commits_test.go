package git_test

import (
	"testing"

	"github.com/Masterminds/semver"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		repository, err := testutil.SetUpInMemoryRepository()
		require.NoError(t, err)

		worktree, err := repository.Worktree()
		require.NoError(t, err)

		for _, commit := range test.commitHistory {
			commitOptions := testutil.CreateCommitOptions()
			_, err := worktree.Commit(commit.message, commitOptions)
			require.NoError(t, err)

			if commit.tag == "" {
				continue
			}

			head, err := repository.Head()
			require.NoError(t, err)

			var createTagOpts *gogit.CreateTagOptions
			if test.annotateTags {
				createTagOpts = &gogit.CreateTagOptions{
					Message: "some message",
					Tagger:  commitOptions.Author,
				}
			}
			_, err = repository.CreateTag(commit.tag, head.Hash(), createTagOpts)
			require.NoError(t, err)
		}

		classifier := conventionalcommits.NewTypeClassifier()
		actual, err := git.GetConventionalCommitTypesSinceLastRelease(repository, classifier)

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
		assert.ElementsMatch(t, test.expectedConventionalCommitTypes, actual.ConventionalCommitTypes)
	}
}

func TestGetConventionalCommitTypesSinceLatestReleaseInShallowRepository(t *testing.T) {
	setUpShallowRepository := func(t *testing.T, commitHistory []commit) *gogit.Repository {
		repository, err := testutil.SetUpInMemoryRepository()
		require.NoError(t, err)

		worktree, err := repository.Worktree()
		require.NoError(t, err)

		for _, currentCommit := range commitHistory {
			_, err := worktree.Commit(currentCommit.message, testutil.CreateCommitOptions())
			require.NoError(t, err)

			if currentCommit.tag == "" {
				continue
			}

			head, err := repository.Head()
			require.NoError(t, err)

			_, err = repository.CreateTag(currentCommit.tag, head.Hash(), nil)
			require.NoError(t, err)
		}

		// Simulate the truncated history of a shallow clone by marking the
		// oldest known commit as a shallow boundary.
		commits, err := repository.Log(&gogit.LogOptions{Order: gogit.LogOrderCommitterTime})
		require.NoError(t, err)

		var oldestCommitHash plumbing.Hash
		require.NoError(t, commits.ForEach(func(currentCommit *object.Commit) error {
			oldestCommitHash = currentCommit.Hash
			return nil
		}))

		require.NoError(t, repository.Storer.SetShallow([]plumbing.Hash{oldestCommitHash}))

		return repository
	}

	classifier := conventionalcommits.NewTypeClassifier()

	t.Run("returns an error if no release tag is within the truncated history", func(t *testing.T) {
		repository := setUpShallowRepository(t, []commit{
			{message: "chore: Do something", tag: ""},
			{message: "feat: Do something else", tag: ""},
		})

		_, err := git.GetConventionalCommitTypesSinceLastRelease(repository, classifier)

		assert.ErrorIs(t, err, git.ErrShallowRepository)
	})

	t.Run("returns the result if a release tag is within the truncated history", func(t *testing.T) {
		repository := setUpShallowRepository(t, []commit{
			{message: "chore: Last release", tag: "1.0.0"},
			{message: "feat: Do something", tag: ""},
		})

		actual, err := git.GetConventionalCommitTypesSinceLastRelease(repository, classifier)

		require.NoError(t, err)
		assert.True(t, semver.MustParse("1.0.0").Equal(actual.LatestReleaseVersion))
		assert.ElementsMatch(t, []conventionalcommits.Type{conventionalcommits.Feature}, actual.ConventionalCommitTypes)
	})

	t.Run("does not report a complete repository without any release tags as shallow", func(t *testing.T) {
		repository, err := testutil.SetUpInMemoryRepository()
		require.NoError(t, err)

		worktree, err := repository.Worktree()
		require.NoError(t, err)

		_, err = worktree.Commit("feat: Do something", testutil.CreateCommitOptions())
		require.NoError(t, err)

		actual, err := git.GetConventionalCommitTypesSinceLastRelease(repository, classifier)

		require.NoError(t, err)
		assert.True(t, semver.MustParse("0.0.0").Equal(actual.LatestReleaseVersion))
	})
}
