package conventionalcommits_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
)

func TestStringToType(t *testing.T) {
	classifier := conventionalcommits.NewTypeClassifier()
	
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
		commitType, err := classifier.StringToType(test.string)

		if test.doExpectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, test.expectedType, commitType)
	}
}

func TestTypeClassifier(t *testing.T) {
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
			customChoreTypes:   []string{"build", "chore", "ci", "docs", "style", "refactor", "test"},
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
			testString:         "fix",
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
		{
			name:               "precedence: fix takes precedence over chore when in both",
			customFixTypes:     []string{"fix", "perf"},
			testString:         "perf",
			expectedType:       conventionalcommits.Fix,
			doExpectError:      false,
		},
		{
			name:               "precedence: feature takes precedence over fix when in both",
			customFixTypes:     []string{"fix", "enhance"},
			customFeatureTypes: []string{"feat", "enhance"},
			testString:         "enhance",
			expectedType:       conventionalcommits.Feature,
			doExpectError:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			classifier := conventionalcommits.NewTypeClassifierWithCustomPrefixes(test.customChoreTypes, test.customFixTypes, test.customFeatureTypes)
			
			commitType, err := classifier.StringToType(test.testString)
			
			if test.doExpectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedType, commitType)
			}
		})
	}
}

func TestNewTypeClassifier(t *testing.T) {
	classifier := conventionalcommits.NewTypeClassifier()
	
	commitType, err := classifier.StringToType("feat")
	assert.NoError(t, err)
	assert.Equal(t, conventionalcommits.Feature, commitType)
	
	commitType, err = classifier.StringToType("fix")
	assert.NoError(t, err)
	assert.Equal(t, conventionalcommits.Fix, commitType)
	
	commitType, err = classifier.StringToType("chore")
	assert.NoError(t, err)
	assert.Equal(t, conventionalcommits.Chore, commitType)
}
