package target_test

import (
	"io"
	"os"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"

	"github.com/thenativeweb/get-next-version/target"
)

func TestWriteOutput(t *testing.T) {
	version, err := semver.NewVersion("1.2.3")
	assert.NoError(t, err)

	t.Run("writes output to the github output file", func(t *testing.T) {
		outputFile, err := os.CreateTemp("", "get-next-version-*")
		if err != nil {
			panic(err)
		}
		githubOutputFile := outputFile.Name()
		outputFile.Close()
		defer os.Remove(outputFile.Name())
		os.Setenv("GITHUB_OUTPUT", githubOutputFile)

		err = target.WriteOutput(*version, true, "github-action", "")
		assert.NoError(t, err)

		outputFile, err = os.Open(githubOutputFile)
		assert.NoError(t, err)
		data, err := io.ReadAll(outputFile)
		assert.NoError(t, err)
		assert.Equal(t, "version=1.2.3\nhasNextVersion=true\n", string(data))
	})

	t.Run("appends output to the github output file", func(t *testing.T) {
		outputFile, err := os.CreateTemp("", "get-next-version-*")
		if err != nil {
			panic(err)
		}
		outputFile.WriteString("prefix foo\n")
		githubOutputFile := outputFile.Name()
		outputFile.Close()
		defer os.Remove(outputFile.Name())
		os.Setenv("GITHUB_OUTPUT", githubOutputFile)

		err = target.WriteOutput(*version, true, "github-action", "")
		assert.NoError(t, err)

		outputFile, err = os.Open(githubOutputFile)
		assert.NoError(t, err)
		data, err := io.ReadAll(outputFile)
		assert.NoError(t, err)
		assert.Equal(t, "prefix foo\nversion=1.2.3\nhasNextVersion=true\n", string(data))
	})

	t.Run("returns an error if the GITHUB_OUTPUT environment variable is not set", func(t *testing.T) {
		os.Setenv("GITHUB_OUTPUT", "")

		err = target.WriteOutput(*version, true, "github-action", "")
		assert.EqualError(t, err, "environment variable GITHUB_OUTPUT must be set")
	})

	t.Run("returns an error if the output file cannot be opened for writing", func(t *testing.T) {
		os.Setenv("GITHUB_OUTPUT", "/root")

		err = target.WriteOutput(*version, true, "github-action", "")
		assert.EqualError(t, err, "could not open github output file for writing: open /root: is a directory")
	})
}
