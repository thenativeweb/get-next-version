package testutil

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func SetUpInMemoryRepository() *git.Repository {
	storer := memory.NewStorage()
	fs := memfs.New()

	repository, _ := git.Init(storer, fs)
	worktree, _ := repository.Worktree()
	worktree.Commit("Initial Commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
	})

	return repository
}
