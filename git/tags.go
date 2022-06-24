package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetAllTags(repository *git.Repository) ([]*plumbing.Reference, error) {
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
