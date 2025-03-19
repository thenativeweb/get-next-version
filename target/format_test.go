package target_test

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/target"
)

func TestFormat(t *testing.T) {
	version, err := semver.NewVersion("1.2.3")
	assert.NoError(t, err)

	output := target.Format(*version, true, "github-action", "")
	assert.Equal(t, []string{
		"version=1.2.3",
		"hasNextVersion=true",
	}, output)

	output = target.Format(*version, false, "github-action", "")
	assert.Equal(t, []string{
		"version=1.2.3",
		"hasNextVersion=false",
	}, output)

	output = target.Format(*version, false, "github-action", "v")
	assert.Equal(t, []string{
		"version=v1.2.3",
		"hasNextVersion=false",
	}, output)

	output = target.Format(*version, true, "json", "")
	assert.Equal(t, []string{
		`{"version": "1.2.3", "hasNextVersion": true}`,
	}, output)

	output = target.Format(*version, false, "json", "")
	assert.Equal(t, []string{
		`{"version": "1.2.3", "hasNextVersion": false}`,
	}, output)

	output = target.Format(*version, false, "json", "v")
	assert.Equal(t, []string{
		`{"version": "v1.2.3", "hasNextVersion": false}`,
	}, output)

	output = target.Format(*version, true, "version", "")
	assert.Equal(t, []string{
		"1.2.3",
	}, output)

	output = target.Format(*version, false, "version", "")
	assert.Equal(t, []string{
		"1.2.3",
	}, output)

	output = target.Format(*version, false, "version", "v")
	assert.Equal(t, []string{
		"v1.2.3",
	}, output)

	assert.Panics(t, func() {
		target.Format(*version, true, "non-existent-format", "")
	})
}
