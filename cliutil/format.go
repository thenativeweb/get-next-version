package cliutil

import (
	"fmt"

	"github.com/Masterminds/semver"
)

func Format(nextVersion semver.Version, hasNextVersion bool, format string) []string {
	switch format {
	case "github-action":
		return []string{
			fmt.Sprintf("::set-output name=version::%s", nextVersion.String()),
			fmt.Sprintf("::set-output name=hasNextVersion::%v", hasNextVersion),
		}
	case "json":
		return []string{
			fmt.Sprintf(`{"version": "%s", "hasNextVersion": %v}`, nextVersion.String(), hasNextVersion),
		}
	case "version":
		return []string{
			nextVersion.String(),
		}
	default:
		panic("invalid format")
	}
}
