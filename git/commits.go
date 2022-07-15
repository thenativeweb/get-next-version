package git

import (
	"errors"
	"fmt"
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
	fmt.Println("[DEBUG]: Get All Tags")
	tags, err := GetAllTags(repository)
	if err != nil {
		fmt.Printf("[DEBUG] Error: %v\n", err)
		return ConventionalCommmitTypesResult{}, err
	}
	fmt.Printf("[DEBUG] Tags: %v\n", tags)

	fmt.Println("[DEBUG]: Getting Head of Repo")
	head, err := repository.Head()
	if err != nil {
		fmt.Printf("[DEBUG] Error: %v\n", err)
		return ConventionalCommmitTypesResult{}, err
	}

	fmt.Println("[DEBUG]: Getting Commit Iterator")
	commitIterator, err := repository.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		fmt.Printf("[DEBUG] Error: %v\n", err)
		return ConventionalCommmitTypesResult{}, err
	}
	fmt.Printf("[DEBUG] Commit Iterator: %v\n", commitIterator)

	fmt.Println("[DEBUG]: Getting Next Commit")
	currentCommit, currentCommitErr := commitIterator.Next()
	var latestReleaseVersion *semver.Version
	conventionalCommitTypes := []conventionalcommits.Type{}

	for currentCommitErr == nil {
		fmt.Printf("[DEBUG] Current Commit: %v\n", currentCommit)
		wasPartOfLastRelease := false
		for _, tag := range tags {
			if tag.Hash() == currentCommit.Hash {
				fmt.Println("[DEBUG]: Found Tag associated with Hash")
				fmt.Printf("[DEBUG] Tag: %v\n", tag)

				fmt.Println("[DEBUG]: Checking if Tag is Version")
				latestReleaseVersion, err = semver.NewVersion(tag.Name().Short())
				if err == nil {
					fmt.Println("[DEBUG]: Detected Current Version")
					fmt.Printf("[DEBUG] Release Version: %v\n", latestReleaseVersion)
					wasPartOfLastRelease = true
					break
				}
			}
		}

		if wasPartOfLastRelease {
			fmt.Println("[DEBUG]: Getting Next Commit")
			break
		}

		fmt.Println("[DEBUG]: Getting Commit Type for Commit")
		currentCommitType, err := conventionalcommits.CommitMessageToType(currentCommit.Message)
		if err != nil {
			currentCommitType = conventionalcommits.Chore
		}
		fmt.Printf("[DEBUG] Commit Type: %v\n", currentCommitType)

		conventionalCommitTypes = append(
			conventionalCommitTypes,
			currentCommitType,
		)
		currentCommit, currentCommitErr = commitIterator.Next()
	}

	if currentCommitErr != nil {
		fmt.Printf("[DEBUG] Commit Error: %v\n", currentCommitErr)
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
