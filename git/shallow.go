package git

import (
	"github.com/go-git/go-git/v5"
)

func IsShallow(repository *git.Repository) (bool, error) {
	shallowCommits, err := repository.Storer.Shallow()
	if err != nil {
		return false, err
	}

	return len(shallowCommits) > 0, nil
}
