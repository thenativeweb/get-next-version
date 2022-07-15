package conventionalcommits_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
)

func TestCommitMessageToType(t *testing.T) {
	tests := []struct {
		message            string
		doExpectError      bool
		expectedCommitType conventionalcommits.Type
	}{
		{message: "chore:", doExpectError: false, expectedCommitType: conventionalcommits.Chore},
		{message: "fix:", doExpectError: false, expectedCommitType: conventionalcommits.Fix},
		{message: "feat:", doExpectError: false, expectedCommitType: conventionalcommits.Feature},
		{message: "chore!:", doExpectError: false, expectedCommitType: conventionalcommits.BreakingChange},
		{message: "chore(scope):", doExpectError: false, expectedCommitType: conventionalcommits.Chore},
		{message: "chore(scope)!:", doExpectError: false, expectedCommitType: conventionalcommits.BreakingChange},
		{message: "chore:\n\nBREAKING CHANGE: ", doExpectError: false, expectedCommitType: conventionalcommits.BreakingChange},
		{message: "chore:\n\nBREAKING CHANGE #", doExpectError: false, expectedCommitType: conventionalcommits.BreakingChange},
		{message: "chore:\n\nBREAKING-CHANGE: ", doExpectError: false, expectedCommitType: conventionalcommits.BreakingChange},
		{message: "chore:\n\nBREAKING-CHANGE", doExpectError: false, expectedCommitType: conventionalcommits.Chore},
		{message: "chore:\n\nBREAKING-CHANGE: \n\nSome-Token: ", doExpectError: false, expectedCommitType: conventionalcommits.BreakingChange},
		{message: "chore:Some Description\nBREAKING-CHANGE: ", doExpectError: false, expectedCommitType: conventionalcommits.Chore},
		{message: "chore:\n\nSome-Token: ", doExpectError: false, expectedCommitType: conventionalcommits.Chore},
		{message: "Chore:", doExpectError: false, expectedCommitType: conventionalcommits.Chore},
		{message: "", doExpectError: true, expectedCommitType: conventionalcommits.Chore},
		{message: "invaild:", doExpectError: true, expectedCommitType: conventionalcommits.Chore},
	}

	for _, test := range tests {
		commitType, err := conventionalcommits.CommitMessageToType(test.message)

		if test.doExpectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, test.expectedCommitType, commitType)
	}
}
