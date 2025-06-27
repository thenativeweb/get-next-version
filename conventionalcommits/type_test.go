package conventionalcommits_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
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

		assert.NoError(t, err)
		assert.Equal(t, test.expectedType, commitType)
	}
}

func TestSetCustomPrefixes(t *testing.T) {
	// Reset to defaults before each test
	conventionalcommits.ResetToDefaults()
	
	tests := []struct {
		name                 string
		customChoreTypes     []string
		customFixTypes       []string  
		customFeatureTypes   []string
		testString           string
		expectedType         conventionalcommits.Type
		doExpectError        bool
	}{
		{
			name:               "custom fix prefix deps",
			customFixTypes:     []string{"fix", "deps"},
			testString:         "deps",
			expectedType:       conventionalcommits.Fix,
			doExpectError:      false,
		},
		{
			name:               "custom fix prefix perf moves from chore to fix", 
			customFixTypes:     []string{"fix", "perf"},
			customChoreTypes:   []string{"build", "chore", "ci", "docs", "style", "refactor", "test"}, // Remove perf from chore
			testString:         "perf",
			expectedType:       conventionalcommits.Fix,
			doExpectError:      false,
		},
		{
			name:               "custom feature prefix enhance",
			customFeatureTypes: []string{"feat", "enhance"},
			testString:         "enhance", 
			expectedType:       conventionalcommits.Feature,
			doExpectError:      false,
		},
		{
			name:               "custom chore prefix update",
			customChoreTypes:   []string{"chore", "update"},
			testString:         "update",
			expectedType:       conventionalcommits.Chore,
			doExpectError:      false,
		},
		{
			name:               "override defaults completely",
			customFixTypes:     []string{"patch"},
			customFeatureTypes: []string{"minor"},
			testString:         "fix", // should now be invalid since we overrode
			expectedType:       conventionalcommits.Chore,
			doExpectError:      true,
		},
		{
			name:               "override defaults - new prefix works",
			customFixTypes:     []string{"patch"},
			testString:         "patch",
			expectedType:       conventionalcommits.Fix,
			doExpectError:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset before each subtest
			conventionalcommits.ResetToDefaults()
			
			// Set custom prefixes
			conventionalcommits.SetCustomPrefixes(test.customChoreTypes, test.customFixTypes, test.customFeatureTypes)
			
			// Test the string conversion
			commitType, err := conventionalcommits.StringToType(test.testString)
			
			if test.doExpectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedType, commitType)
			}
		})
	}
}

func TestResetToDefaults(t *testing.T) {
	// Set custom prefixes
	conventionalcommits.SetCustomPrefixes([]string{"custom"}, []string{"custom"}, []string{"custom"})
	
	// Reset to defaults
	conventionalcommits.ResetToDefaults()
	
	// Test that default prefixes work again
	commitType, err := conventionalcommits.StringToType("feat")
	assert.NoError(t, err)
	assert.Equal(t, conventionalcommits.Feature, commitType)
	
	commitType, err = conventionalcommits.StringToType("fix")
	assert.NoError(t, err)
	assert.Equal(t, conventionalcommits.Fix, commitType)
	
	commitType, err = conventionalcommits.StringToType("chore")
	assert.NoError(t, err)
	assert.Equal(t, conventionalcommits.Chore, commitType)
}
