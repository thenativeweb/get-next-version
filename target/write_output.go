package target

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Masterminds/semver"
)

func WriteOutput(nextVersion semver.Version, hasNextVersion bool, target string, prefix string) error {
	outputLines := Format(nextVersion, hasNextVersion, target, prefix)

	var outputHandle *os.File
	switch target {
	case "github-action":
		githubOutputFile := os.Getenv("GITHUB_OUTPUT")
		if githubOutputFile == "" {
			return fmt.Errorf("environment variable GITHUB_OUTPUT must be set")
		}

		var err error
		outputHandle, err = os.OpenFile(githubOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("could not open github output file for writing: %w", err)
		}
		defer outputHandle.Close()
	case "json":
		outputHandle = os.Stdout
	case "version":
		outputHandle = os.Stdout
	default:
		panic("invalid target")
	}

	writer := bufio.NewWriter(outputHandle)
	for _, line := range outputLines {
		if _, err := writer.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return fmt.Errorf("could not write to target: %w", err)
		}
	}
	writer.Flush()

	return nil
}
