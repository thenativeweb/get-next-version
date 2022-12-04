package git

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Tags = map[plumbing.Hash]*semver.Version

func GetAllSemVerTags(repository *git.Repository) (Tags, error) {
	tagsIterator, err := repository.Tags()
	if err != nil {
		return Tags{}, err
	}

	var tags = make(Tags)

	err = tagsIterator.ForEach(func(tag *plumbing.Reference) error {
		var commitHash plumbing.Hash
		tagObject, err := repository.TagObject(tag.Hash())
		switch err {
		case nil:
			commit, err := tagObject.Commit()
			if err != nil {
				return err
			}
			commitHash = commit.Hash
		case plumbing.ErrObjectNotFound:
			commitHash = tag.Hash()
		default:
			return err
		}

		version, err := semver.NewVersion(tag.Name().Short())
		if err != nil {
			return nil
		}

		if _, exists := tags[commitHash]; exists {
			return errors.New(fmt.Sprintf("commit %s was tagged with multiple semver versions", commitHash.String()))
		}
		tags[commitHash] = version
		return nil
	})
	if err != nil {
		return Tags{}, err
	}

	return tags, nil
}

/*
IsValidTagName runs a series of regex checks to ensure that the tag name is valid
  Tags cannot begin or end with, or contain multiple consecutive / characters.
  They cannot contain any of the following characters \, ?, ~, ^, :, * , [, @.
  They cannot contain a space.
  They cannot end with a . or have two consecutive .. anywhere within them.
  Tags are not case-sensitive.
*/
func IsValidTagName(tag string) bool {
	var mustNotContainBlacklistedChars = regexp.MustCompile(`\?|\\|\~|\^|\:|\*|\[|\@|\s`)
	var mustNotEndWithDot = regexp.MustCompile(`\.$`)
	var mustNotIncludeDoubleSlash = regexp.MustCompile(`\/{2,}`)
	var mustNotIncludeDoubleDots = regexp.MustCompile(`\.{2,}`)
	var mustNotStartOrEndWithSlash = regexp.MustCompile(`^\/|\/$`)

	return !mustNotContainBlacklistedChars.MatchString(tag) &&
		!mustNotEndWithDot.MatchString(tag) &&
		!mustNotIncludeDoubleSlash.MatchString(tag) &&
		!mustNotIncludeDoubleDots.MatchString(tag) &&
		!mustNotStartOrEndWithSlash.MatchString(tag)
}
