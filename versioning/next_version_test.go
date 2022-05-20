package versioning_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/getnextversion/versioning"
)

func TestCalculateNextVersion(t *testing.T) {
	tests := []struct {
		currentVersion         *semver.Version
		semanticVersioningTags []versioning.SemanticVersioningTag
		expectedNewVersion     *semver.Version
	}{
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.1"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.1.1"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Chore, versioning.Fix, versioning.Feature, versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.BreakingChange, versioning.Feature, versioning.Fix, versioning.Chore},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Chore, versioning.Chore},
			expectedNewVersion:     semver.MustParse("1.0.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Fix, versioning.Fix},
			expectedNewVersion:     semver.MustParse("1.0.1"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.Feature, versioning.Feature},
			expectedNewVersion:     semver.MustParse("1.1.0"),
		},
		{
			currentVersion:         semver.MustParse("1.0.0"),
			semanticVersioningTags: []versioning.SemanticVersioningTag{versioning.BreakingChange, versioning.BreakingChange},
			expectedNewVersion:     semver.MustParse("2.0.0"),
		},
	}

	for _, test := range tests {

		actualNewVersion := versioning.CalculateNextVersion(test.currentVersion, test.semanticVersioningTags)
		assert.Equal(t, test.expectedNewVersion, &actualNewVersion)
	}

}
