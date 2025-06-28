package git

import (
	"errors"
	"fmt"
	"strings"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Tags = map[plumbing.Hash]*semver.Version

// tagCandidate represents a potential semver tag for a commit.
type tagCandidate struct {
	originalName string
	version      *semver.Version
}

// getTagSpecificity returns the number of version components in a tag.
func getTagSpecificity(tagName string) int {
	cleanTag := tagName
	if strings.HasPrefix(cleanTag, "v") {
		cleanTag = cleanTag[1:]
	}
	return strings.Count(cleanTag, ".")
}

// areCompatibleGranularities determines if two versions represent different granularity levels
// of the same logical version rather than conflicting versions.
func areCompatibleGranularities(leftVersion, rightVersion *semver.Version) bool {
	if leftVersion.Equal(rightVersion) {
		return true
	}
	
	if leftVersion.Major() != rightVersion.Major() {
		return false
	}
	
	// Allow major-only tags (e.g., "v4" represented as v4.0.0) 
	if (leftVersion.Minor() == 0 && leftVersion.Patch() == 0) || (rightVersion.Minor() == 0 && rightVersion.Patch() == 0) {
		return true
	}
	
	if leftVersion.Minor() == rightVersion.Minor() {
		// Allow major.minor tags (e.g., "v4.5" represented as v4.5.0)
		if leftVersion.Patch() == 0 || rightVersion.Patch() == 0 {
			return true
		}
		return leftVersion.Patch() == rightVersion.Patch()
	}
	
	return false
}

// selectMostSpecificTag selects the tag with the highest specificity (most version components).
func selectMostSpecificTag(candidates []tagCandidate) *semver.Version {
	if len(candidates) == 1 {
		return candidates[0].version
	}
	
	mostSpecific := candidates[0]
	maxSpecificity := getTagSpecificity(mostSpecific.originalName)
	
	for _, candidate := range candidates[1:] {
		specificity := getTagSpecificity(candidate.originalName)
		if specificity > maxSpecificity {
			mostSpecific = candidate
			maxSpecificity = specificity
		}
	}
	
	return mostSpecific.version
}

// GetAllSemVerTags extracts semantic version tags from a repository.
//
// Algorithm: When multiple tags exist on the same commit, this function distinguishes
// between acceptable granularity variations (e.g., v4, v4.5, v4.5.14) and conflicting 
// versions (e.g., v4.1.0, v4.2.0). For granularity variations, it selects the most
// specific tag. For conflicting versions, it returns an error.
func GetAllSemVerTags(repository *git.Repository) (Tags, error) {
	tagsIterator, err := repository.Tags()
	if err != nil {
		return Tags{}, err
	}

	var commitTags = make(map[plumbing.Hash][]tagCandidate)

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
			return nil // Skip non-semver tags
		}

		commitTags[commitHash] = append(commitTags[commitHash], tagCandidate{
			originalName: tag.Name().Short(),
			version:      version,
		})
		return nil
	})
	if err != nil {
		return Tags{}, err
	}

	var tags = make(Tags)
	for commitHash, candidates := range commitTags {
		if len(candidates) > 1 {
			firstVersion := candidates[0].version
			hasDifferentVersions := false
			
			for _, candidate := range candidates[1:] {
				if !areCompatibleGranularities(firstVersion, candidate.version) {
					hasDifferentVersions = true
					break
				}
			}
			
			if hasDifferentVersions {
				return Tags{}, errors.New(fmt.Sprintf("commit %s was tagged with multiple semver versions", commitHash.String()))
			}
		}
		
		tags[commitHash] = selectMostSpecificTag(candidates)
	}

	return tags, nil
}
