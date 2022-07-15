package git_test

import (
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/testutil"
)

func TestGetAllTags(t *testing.T) {
	tests := []struct {
		tagsPerBranch    map[string][]string
		expectedTagNames []string
	}{
		{
			tagsPerBranch:    map[string][]string{"main": {"1.0.0"}},
			expectedTagNames: []string{"1.0.0"},
		},
		{
			tagsPerBranch:    map[string][]string{"main": {"1.0.0", "2.0.0"}},
			expectedTagNames: []string{"1.0.0", "2.0.0"},
		},
		{
			tagsPerBranch:    map[string][]string{"main": {"1.0.0", "2.0.0"}, "feature": {"feature-tag"}},
			expectedTagNames: []string{"1.0.0", "2.0.0", "feature-tag"},
		},
	}

	for _, test := range tests {
		repository := testutil.SetUpInMemoryRepository()

		for branchName, tagNames := range test.tagsPerBranch {
			worktree, _ := repository.Worktree()
			worktree.Checkout(&gogit.CheckoutOptions{
				Create: true,
				Branch: plumbing.ReferenceName(branchName),
			})

			head, _ := repository.Head()
			for _, tagName := range tagNames {
				repository.CreateTag(tagName, head.Hash(), nil)
			}
		}

		tags, err := git.GetAllTags(repository)

		assert.NoError(t, err)
		var tagNames []string
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name().Short())
		}
		assert.ElementsMatch(t, test.expectedTagNames, tagNames)
	}
}
