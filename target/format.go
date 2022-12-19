package target

import (
	"fmt"
	"github.com/Masterminds/semver"
)

func Format(nextVersion semver.Version, hasNextVersion bool, format string, prefix string) []string {
	versionString := prefix + nextVersion.String()

	switch format {
	case "github-action":
		return []string{
			fmt.Sprintf("version=%s", versionString),
			fmt.Sprintf("hasNextVersion=%v", hasNextVersion),
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
