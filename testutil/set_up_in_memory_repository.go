package testutil

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func SetUpInMemoryRepository() (*git.Repository, error) {
	storer := memory.NewStorage()
	fs := memfs.New()

	repository, err := git.Init(storer, fs)
	if err != nil {
		return nil, err
	}

	return repository, nil
}
