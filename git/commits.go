package git

import (
	"errors"
	"io"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/thenativeweb/getnextversion/conventionalcommits"
)

type ConventionalCommmitTypesResult struct {
	LatestReleaseVersion    *semver.Version
	ConventionalCommitTypes []conventionalcommits.Type
}

var ErrNoCommitsFound = errors.New("no commits found")

func GetConventionalCommitTypesSinceLastRelease(repository *git.Repository) (ConventionalCommmitTypesResult, error) {
	tags, err := GetAllTags(repository)
	if err != nil {
		return ConventionalCommmitTypesResult{}, err
	}

	head, err := repository.Head()
	if err != nil {
		return ConventionalCommmitTypesResult{}, err
	}

	commitIterator, err := repository.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return ConventionalCommmitTypesResult{}, err
	}

	currentCommit, currentCommitErr := commitIterator.Next()
	var latestReleaseVersion *semver.Version
	conventionalCommitTypes := []conventionalcommits.Type{}

	for currentCommitErr == nil {
		wasPartOfLastRelease := false
		for _, tag := range tags {
			if tag.Hash() == currentCommit.Hash {

				latestReleaseVersion, err = semver.NewVersion(tag.Name().Short())
				if err == nil {
					wasPartOfLastRelease = true
					break
				}
			}
		}

		if wasPartOfLastRelease {
			break
		}

		currentCommitType, err := conventionalcommits.CommitMessageToType(currentCommit.Message)
		if err != nil {
			currentCommitType = conventionalcommits.Chore
		}

		conventionalCommitTypes = append(
			conventionalCommitTypes,
			currentCommitType,
		)
		currentCommit, currentCommitErr = commitIterator.Next()
	}

	if currentCommitErr != nil {
		if currentCommitErr == io.EOF && currentCommit == nil {
			return ConventionalCommmitTypesResult{}, ErrNoCommitsFound
		}

		return ConventionalCommmitTypesResult{}, currentCommitErr
	}

	return ConventionalCommmitTypesResult{
		LatestReleaseVersion:    latestReleaseVersion,
		ConventionalCommitTypes: conventionalCommitTypes,
	}, nil
}
