package cliutil

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
)

func Format(nextVersion semver.Version, hasNextVersion bool, format string, prefix ...string) []string {
	prefixString := strings.Join(prefix, "")
	versionString := prefixString + nextVersion.String()

	switch format {
	case "github-action":
		return []string{
			fmt.Sprintf("::set-output name=version::%s", versionString),
			fmt.Sprintf("::set-output name=hasNextVersion::%v", hasNextVersion),
		}
	case "json":
		return []string{
			fmt.Sprintf(`{"version": "%s", "hasNextVersion": %v}`, versionString, hasNextVersion),
		}
	case "version":
		return []string{
			versionString,
		}
	default:
		panic("invalid format")
	}
}
