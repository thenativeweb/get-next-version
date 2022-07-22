package testutil

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func SetUpInMemoryRepository() *git.Repository {
	storer := memory.NewStorage()
	fs := memfs.New()

	repository, _ := git.Init(storer, fs)
	return repository
}
