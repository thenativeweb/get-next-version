package cliutil

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/thenativeweb/get-next-version/git"
)

func Format(nextVersion semver.Version, hasNextVersion bool, changelog []string, format string) []string {
	changelogJson, _ := json.Marshal(changelog)

	switch format {
	case "github-action":
		return []string{
			fmt.Sprintf("::set-output name=version::%s", nextVersion.String()),
			fmt.Sprintf("::set-output name=hasNextVersion::%v", hasNextVersion),
			fmt.Sprintf("::set-output name=changelog::%s", string(changelogJson)),
		}
	case "json":
		return []string{
			fmt.Sprintf(`{"version": "%s", "hasNextVersion": %v, "changelog": %s}`, nextVersion.String(), hasNextVersion, string(changelogJson)),
		}
	case "version":
		return []string{
			nextVersion.String(),
		}
	default:
		panic("invalid format")
	}
}

func ExtractChangelog(commits []git.ConventionalCommit) []string {
	changelogs := []string{}
	for _, commit := range commits {
		changelogs = append(changelogs, commit.Message)
	}
	return changelogs
}
