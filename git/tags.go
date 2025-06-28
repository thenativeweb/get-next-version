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

// tagCandidate holds a tag candidate with its original name and parsed version
type tagCandidate struct {
	originalName string
	version      *semver.Version
}

// getTagSpecificity returns how specific a tag is based on the number of version components.
// More dots means more specific: "v4" (0 dots) < "v4.5" (1 dot) < "v4.5.14" (2 dots)
func getTagSpecificity(tagName string) int {
	// Remove 'v' prefix if present
	cleanTag := tagName
	if strings.HasPrefix(cleanTag, "v") {
		cleanTag = cleanTag[1:]
	}
	return strings.Count(cleanTag, ".")
}

// areCompatibleGranularities checks if two versions are compatible granularities of the same version.
// For example, v4.0.0, v4.5.0, and v4.5.14 are compatible if they represent granularities like v4, v4.5, v4.5.14
func areCompatibleGranularities(v1, v2 *semver.Version) bool {
	// If both versions are identical, they're definitely compatible
	if v1.Equal(v2) {
		return true
	}
	
	// Check if one could be a less specific version of the other
	// v4.0.0 is compatible with v4.5.14 (v4 vs v4.5.14)
	// v4.5.0 is compatible with v4.5.14 (v4.5 vs v4.5.14)
	// v4.0.0 is NOT compatible with v5.0.0 (different major)
	// v4.1.0 is NOT compatible with v4.2.14 (different minor, both explicit)
	
	if v1.Major() != v2.Major() {
		return false
	}
	
	// If one has minor = 0 and patch = 0, it could be a "major only" tag like "v4"
	if (v1.Minor() == 0 && v1.Patch() == 0) || (v2.Minor() == 0 && v2.Patch() == 0) {
		return true
	}
	
	// If they have the same major and minor, check patch
	if v1.Minor() == v2.Minor() {
		// If one has patch = 0, it could be a "major.minor" tag like "v4.5"
		if v1.Patch() == 0 || v2.Patch() == 0 {
			return true
		}
		// If both have explicit patch versions, they should be the same
		return v1.Patch() == v2.Patch()
	}
	
	// Different minor versions with explicit values are not compatible
	return false
}

// selectMostSpecificTag selects the most specific tag from a list of candidates
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

func GetAllSemVerTags(repository *git.Repository) (Tags, error) {
	tagsIterator, err := repository.Tags()
	if err != nil {
		return Tags{}, err
	}

	// First pass: collect all tag candidates for each commit
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

		// Add this tag as a candidate for this commit
		commitTags[commitHash] = append(commitTags[commitHash], tagCandidate{
			originalName: tag.Name().Short(),
			version:      version,
		})
		return nil
	})
	if err != nil {
		return Tags{}, err
	}

	// Second pass: select the most specific tag for each commit
	var tags = make(Tags)
	for commitHash, candidates := range commitTags {
		if len(candidates) > 1 {
			// Check if this is the problematic case where we have multiple different semver versions
			// vs. the acceptable case where we have multiple granularity tags for the same version
			firstVersion := candidates[0].version
			hasDifferentVersions := false
			
			for _, candidate := range candidates[1:] {
				// If the major, minor, patch don't match in a way that suggests they're different versions
				// (not just different granularities of the same version), then it's an error
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
