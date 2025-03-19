package testutil

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CreateCommitOptions() *git.CommitOptions {
	return &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		AllowEmptyCommits: true,
	}
}
