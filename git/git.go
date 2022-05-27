package git

import (
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/thenativeweb/getnextversion/versioning"
)

type ConventionalCommmitTypesResult struct {
	LatestReleaseVersion    *semver.Version
	ConventionalCommitTypes []versioning.ConventionalCommitType
}

func getAllTags(repository *git.Repository) ([]*plumbing.Reference, error) {
	tagsIterator, err := repository.Tags()
	if err != nil {
		return []*plumbing.Reference{}, err
	}

	var tags []*plumbing.Reference

	tagsIterator.ForEach(func(tag *plumbing.Reference) error {
		tags = append(tags, tag)
		return nil
	})

	return tags, nil
}

func GetConventionalCommitTypesSinceLastRelease(repository *git.Repository) (ConventionalCommmitTypesResult, error) {
	tags, err := getAllTags(repository)
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
	var conventionalCommitTypes []versioning.ConventionalCommitType
	for currentCommitErr != nil {
		var wasPartOfLastRelease = false
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
		conventionalCommitTypes = append(conventionalCommitTypes, commitMessageToConventionalCommitType(currentCommit.Message))
	}

	if currentCommitErr != nil {
		return ConventionalCommmitTypesResult{}, currentCommitErr
	}

	return ConventionalCommmitTypesResult{
		LatestReleaseVersion:    latestReleaseVersion,
		ConventionalCommitTypes: conventionalCommitTypes,
	}, nil
}

func commitMessageToConventionalCommitType(message string) versioning.ConventionalCommitType {
	// TODO: Implement properly
	return versioning.Chore
}