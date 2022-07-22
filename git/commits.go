package git

import (
	"errors"
	"io"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
)

type ConventionalCommmitTypesResult struct {
	LatestReleaseVersion    *semver.Version
	ConventionalCommitTypes []conventionalcommits.Type
}

var ErrNoCommitsFound = errors.New("no commits found")

func GetConventionalCommitTypesSinceLastRelease(repository *git.Repository) (ConventionalCommmitTypesResult, error) {
	tags, err := GetAllSemVerTags(repository)
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

	foundAtLeastOneCommit := false
	currentCommit, currentCommitErr := commitIterator.Next()
	var latestReleaseVersion *semver.Version
	conventionalCommitTypes := []conventionalcommits.Type{}
	for currentCommitErr == nil {
		foundAtLeastOneCommit = true
		var doesVersionExistForCommit bool
		latestReleaseVersion, doesVersionExistForCommit = tags[currentCommit.Hash]
		if doesVersionExistForCommit {
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
		if currentCommitErr == io.EOF && !foundAtLeastOneCommit {
			return ConventionalCommmitTypesResult{}, ErrNoCommitsFound
		}

		return ConventionalCommmitTypesResult{}, currentCommitErr
	}

	return ConventionalCommmitTypesResult{
		LatestReleaseVersion:    latestReleaseVersion,
		ConventionalCommitTypes: conventionalCommitTypes,
	}, nil
}
