package git_test

import (
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/testutil"
)

func TestGetAllSemVerTags(t *testing.T) {
	tests := []struct {
		tagsPerBranch    map[string][][]string
		doesExpectError  bool
		expectedTagNames []string
	}{
		{
			tagsPerBranch:    map[string][][]string{"main": {{"1.0.0"}}},
			doesExpectError:  false,
			expectedTagNames: []string{"1.0.0"},
		},
		{
			tagsPerBranch:    map[string][][]string{"main": {{"v1.0.0"}}},
			doesExpectError:  false,
			expectedTagNames: []string{"1.0.0"},
		},
		{
			tagsPerBranch:    map[string][][]string{"main": {{"1.0.0"}, {"2.0.0"}}},
			doesExpectError:  false,
			expectedTagNames: []string{"1.0.0", "2.0.0"},
		},
		{
			tagsPerBranch:    map[string][][]string{"main": {{"1.0.0"}, {"2.0.0"}}, "feature": {{"3.0.0"}}},
			doesExpectError:  false,
			expectedTagNames: []string{"1.0.0", "2.0.0", "3.0.0"},
		},
		{
			tagsPerBranch:    map[string][][]string{"main": {{"1.0.0"}, {"2.0.0"}, {"feature-tag"}}},
			doesExpectError:  false,
			expectedTagNames: []string{"1.0.0", "2.0.0"},
		},
		{
			tagsPerBranch:    map[string][][]string{"main": {{"1.0.0"}, {"v2.0.0"}, {"vsomething-else"}}},
			doesExpectError:  false,
			expectedTagNames: []string{"1.0.0", "2.0.0"},
		},
		{
			tagsPerBranch:    map[string][][]string{"main": {{"1.0.0", "2.0.0"}}},
			doesExpectError:  true,
			expectedTagNames: []string{},
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

			for _, tagNamesForCommit := range tagNames {
				worktree.Commit("some message", testutil.CreateCommitOptions())
				head, _ := repository.Head()
				for _, tagName := range tagNamesForCommit {
					repository.CreateTag(tagName, head.Hash(), nil)
				}
			}
		}

		tags, err := git.GetAllSemVerTags(repository, "")

		if test.doesExpectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		var tagNames []string
		for _, tag := range tags {
			tagNames = append(tagNames, tag.String())
		}
		assert.ElementsMatch(t, test.expectedTagNames, tagNames)
	}
}
