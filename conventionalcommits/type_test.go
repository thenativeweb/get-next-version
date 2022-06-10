package conventionalcommits_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/getnextversion/conventionalcommits"
)

func TestStringToType(t *testing.T) {
	tests := []struct {
		string        string
		doExpectError bool
		expectedType  conventionalcommits.Type
	}{
		{string: "", doExpectError: true, expectedType: conventionalcommits.Chore},
		{string: "invalid", doExpectError: true, expectedType: conventionalcommits.Chore},
		{string: "chore", doExpectError: false, expectedType: conventionalcommits.Chore},
		{string: "fix", doExpectError: false, expectedType: conventionalcommits.Fix},
		{string: "feat", doExpectError: false, expectedType: conventionalcommits.Feature},
		{string: "Chore", doExpectError: false, expectedType: conventionalcommits.Chore},
	}

	for _, test := range tests {
		commitType, err := conventionalcommits.StringToType(test.string)

		if test.doExpectError {
			assert.Error(t, err)
			continue
		}

		assert.Nil(t, err)
		assert.Equal(t, test.expectedType, commitType)
	}
}
