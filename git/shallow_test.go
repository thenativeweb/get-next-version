package git_test

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/testutil"
)

func TestIsShallow(t *testing.T) {
	t.Run("returns false for a repository with a complete history", func(t *testing.T) {
		repository, err := testutil.SetUpInMemoryRepository()
		require.NoError(t, err)

		isShallow, err := git.IsShallow(repository)
		assert.NoError(t, err)
		assert.False(t, isShallow)
	})

	t.Run("returns true for a repository with a truncated history", func(t *testing.T) {
		repository, err := testutil.SetUpInMemoryRepository()
		require.NoError(t, err)

		worktree, err := repository.Worktree()
		require.NoError(t, err)

		_, err = worktree.Commit("chore: Do something", testutil.CreateCommitOptions())
		require.NoError(t, err)

		head, err := repository.Head()
		require.NoError(t, err)

		err = repository.Storer.SetShallow([]plumbing.Hash{head.Hash()})
		require.NoError(t, err)

		isShallow, err := git.IsShallow(repository)
		assert.NoError(t, err)
		assert.True(t, isShallow)
	})
}
