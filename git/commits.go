package git

import (
	"errors"
	"io"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
)

type ConventionalCommitTypesResult struct {
	LatestReleaseVersion    *semver.Version
	ConventionalCommitTypes []conventionalcommits.Type
}

var ErrNoCommitsFound = errors.New("no commits found")

func GetConventionalCommitTypesSinceLastRelease(repository *git.Repository, classifier *conventionalcommits.TypeClassifier) (ConventionalCommitTypesResult, error) {
	tags, err := GetAllSemVerTags(repository)
	if err != nil {
		return ConventionalCommitTypesResult{}, err
	}
	head, err := repository.Head()
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return ConventionalCommitTypesResult{}, ErrNoCommitsFound
		}
		return ConventionalCommitTypesResult{}, err
	}
	commitIterator, err := repository.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return ConventionalCommitTypesResult{}, err
	}

	currentCommit, currentCommitErr := commitIterator.Next()
	var latestReleaseVersion *semver.Version
	conventionalCommitTypes := []conventionalcommits.Type{}
	for currentCommitErr == nil {
		var doesVersionExistForCommit bool
		latestReleaseVersion, doesVersionExistForCommit = tags[currentCommit.Hash]
		if doesVersionExistForCommit {
			break
		}

		currentCommitType, err := conventionalcommits.CommitMessageToTypeWithClassifier(currentCommit.Message, classifier)
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
		if currentCommitErr != io.EOF {
			return ConventionalCommitTypesResult{}, currentCommitErr
		}

		latestReleaseVersion = semver.MustParse("0.0.0")
	}

	return ConventionalCommitTypesResult{
		LatestReleaseVersion:    latestReleaseVersion,
		ConventionalCommitTypes: conventionalCommitTypes,
	}, nil
}
