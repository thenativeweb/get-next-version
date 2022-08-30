package cliutil_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/cliutil"
)

func TestFormat(t *testing.T) {
	version, err := semver.NewVersion("1.2.3")
	changelog := []string{"chore: not relevant", "feat: add basic authentication", "feat!: enforce authentication", "fix: fix redirect"}
	assert.NoError(t, err)

	output := cliutil.Format(*version, true, changelog, "github-action")
	assert.Equal(t, []string{
		"::set-output name=version::1.2.3",
		"::set-output name=hasNextVersion::true",
		"::set-output name=changelog::[\"chore: not relevant\",\"feat: add basic authentication\",\"feat!: enforce authentication\",\"fix: fix redirect\"]",
	}, output)

	output = cliutil.Format(*version, false, changelog, "github-action")
	assert.Equal(t, []string{
		"::set-output name=version::1.2.3",
		"::set-output name=hasNextVersion::false",
		"::set-output name=changelog::[\"chore: not relevant\",\"feat: add basic authentication\",\"feat!: enforce authentication\",\"fix: fix redirect\"]",
	}, output)

	output = cliutil.Format(*version, true, changelog, "json")
	assert.Equal(t, []string{
		`{"version": "1.2.3", "hasNextVersion": true, "changelog": ["chore: not relevant","feat: add basic authentication","feat!: enforce authentication","fix: fix redirect"]}`,
	}, output)

	output = cliutil.Format(*version, false, changelog, "json")
	assert.Equal(t, []string{
		`{"version": "1.2.3", "hasNextVersion": false, "changelog": ["chore: not relevant","feat: add basic authentication","feat!: enforce authentication","fix: fix redirect"]}`,
	}, output)

	output = cliutil.Format(*version, true, changelog, "version")
	assert.Equal(t, []string{
		"1.2.3",
	}, output)

	output = cliutil.Format(*version, false, changelog, "version")
	assert.Equal(t, []string{
		"1.2.3",
	}, output)

	assert.Panics(t, func() {
		cliutil.Format(*version, true, changelog, "non-existent-format")
	})
}
