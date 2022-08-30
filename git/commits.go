package git

import (
	"errors"
	"io"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
)

type ConventionalCommit struct {
	Type    conventionalcommits.Type
	Message string
}

type ConventionalCommmitsResult struct {
	LatestReleaseVersion *semver.Version
	ConventionalCommits  []ConventionalCommit
}

var ErrNoCommitsFound = errors.New("no commits found")

func GetConventionalCommitsSinceLastRelease(repository *git.Repository) (ConventionalCommmitsResult, error) {
	tags, err := GetAllSemVerTags(repository)
	if err != nil {
		return ConventionalCommmitsResult{}, err
	}
	head, err := repository.Head()
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return ConventionalCommmitsResult{}, ErrNoCommitsFound
		}
		return ConventionalCommmitsResult{}, err
	}
	commitIterator, err := repository.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return ConventionalCommmitsResult{}, err
	}

	currentCommit, currentCommitErr := commitIterator.Next()
	var latestReleaseVersion *semver.Version
	conventionalCommits := []ConventionalCommit{}
	for currentCommitErr == nil {
		var doesVersionExistForCommit bool
		latestReleaseVersion, doesVersionExistForCommit = tags[currentCommit.Hash]
		if doesVersionExistForCommit {
			break
		}

		currentCommitType, err := conventionalcommits.CommitMessageToType(currentCommit.Message)
		if err != nil {
			currentCommitType = conventionalcommits.Chore
		}
		conventionalCommits = append(
			conventionalCommits,
			ConventionalCommit{
				currentCommitType,
				currentCommit.Message,
			},
		)
		currentCommit, currentCommitErr = commitIterator.Next()
	}

	if currentCommitErr != nil {
		if currentCommitErr != io.EOF {
			return ConventionalCommmitsResult{}, currentCommitErr
		}

		latestReleaseVersion = semver.MustParse("0.0.0")
	}

	return ConventionalCommmitsResult{
		LatestReleaseVersion: latestReleaseVersion,
		ConventionalCommits:  conventionalCommits,
	}, nil
}
