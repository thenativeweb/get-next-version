package conventionalcommits

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/thenativeweb/getnextversion/util"
)

var conventionalCommitBodyRegex *regexp.Regexp

func init() {
	conventionalCommitPrefixes := []string{"fix", "feat", "build", "chore", "ci", "docs", "style", "refector", "perf", "test"}
	conventionalCommitPrefixesRegexString := ""
	for _, prefix := range conventionalCommitPrefixes {
		conventionalCommitPrefixesRegexString += prefix + "|"
	}
	conventionalCommitPrefixesRegexString = strings.TrimSuffix(conventionalCommitPrefixesRegexString, "|")
	conventionalCommitBodyRegexString := fmt.Sprintf("(?P<type>%s)(\\(.*\\))?(?P<breaking>\\!)?:.*", conventionalCommitPrefixesRegexString)

	conventionalCommitBodyRegex = regexp.MustCompile(conventionalCommitBodyRegexString)
}

func CommitMessageToConventionalCommitType(message string) ConventionalCommitType {
	var body string
	var footers []string

	segments := strings.Split(message, "\n")
	var bodySegments []string
	var lastBodyIndex int
	for i, currentSegment := range segments {
		if currentSegment == "" {
			break
		}
		bodySegments = append(bodySegments, currentSegment)
		lastBodyIndex = i
	}
	body = strings.Join(bodySegments, "\n")

	var currentFooterSegments []string
	for _, currentSegment := range segments[lastBodyIndex+1:] {
		if currentSegment == "" {
			footers = append(footers, strings.Join(currentFooterSegments, "\n"))
			currentFooterSegments = []string{}
		}
		currentFooterSegments = append(currentFooterSegments, currentSegment)
	}

	for _, footer := range footers {
		if util.IsOnePrefix(footer, []string{"BREAKING CHANGE: ", "BREAKING CHANGE #"}).IsOnePrefix {
			return BreakingChange
		}
	}

	parsedMesageBody := conventionalCommitBodyRegex.FindStringSubmatch(body)
	if parsedMesageBody == nil {
		return Chore
	}

	breakingIndicator := parsedMesageBody[3]
	if breakingIndicator == "!" {
		return BreakingChange
	}

	return fromString(parsedMesageBody[1])
}
