package git

import (
	"github.com/go-git/go-git/v5"
)

func ReadLatestCommitMessage() (string, error) {
	repository, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}
	ref, err := repository.Head()
	if err != nil {
		return "", err
	}
	cIter, err := repository.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return "", err
	}
	commit, err := cIter.Next()
	if err != nil {
		return "", err
	}

	return commit.Message, nil
}
