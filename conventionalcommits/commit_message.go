package conventionalcommits

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/thenativeweb/get-next-version/util"
)

var (
	bodyRegex             *regexp.Regexp
	breakingFooterTokens  = []string{"BREAKING CHANGE", "BREAKING-CHANGE"}
	footerTokenSeparators = []string{": ", " #"}
)

func initCommitMessage() {
	typesRegexString := ""
	for _, prefix := range defaultClassifier.GetAllTypes() {
		typesRegexString += prefix + "|"
	}
	typesRegexString = strings.TrimSuffix(typesRegexString, "|")
	conventionalCommitBodyRegexString := fmt.Sprintf(
		"(?P<type>(?i:%s))(\\(.*\\))?(?P<breaking>\\!)?:.*",
		typesRegexString,
	)

	bodyRegex = regexp.MustCompile(conventionalCommitBodyRegexString)
}

func splitCommitMessage(message string) (body string, footers []string) {
	segments := strings.Split(message, "\n")
	var bodySegments []string
	var lastBodyIndex int
	for i, currentSegment := range segments {
		lastBodyIndex = i
		if currentSegment == "" {
			break
		}
		bodySegments = append(bodySegments, currentSegment)
	}
	body = strings.Join(bodySegments, "\n")

	var currentFooterSegments []string
	for _, currentSegment := range segments[lastBodyIndex+1:] {
		if currentSegment == "" {
			footers = append(footers, strings.Join(currentFooterSegments, "\n"))
			currentFooterSegments = []string{}
			continue
		}

		currentFooterSegments = append(currentFooterSegments, currentSegment)
	}
	footers = append(footers, strings.Join(currentFooterSegments, "\n"))

	return body, footers
}

func CommitMessageToType(message string) (Type, error) {
	return CommitMessageToTypeWithClassifier(message, defaultClassifier)
}

func CommitMessageToTypeWithClassifier(message string, classifier *TypeClassifier) (Type, error) {
	body, footers := splitCommitMessage(message)

	var breakingFooterPrefixes []string
	for _, token := range breakingFooterTokens {
		for _, separator := range footerTokenSeparators {
			breakingFooterPrefixes = append(breakingFooterPrefixes, token+separator)
		}
	}
	for _, footer := range footers {
		if util.IsOnePrefix(footer, breakingFooterPrefixes).IsOnePrefix {
			return BreakingChange, nil
		}
	}

	parsedMessageBody := bodyRegex.FindStringSubmatch(body)
	if parsedMessageBody == nil {
		return Chore, errors.New("invalid message body for conventional commit message")
	}

	breakingIndicatorIndex := util.MustFind(bodyRegex.SubexpNames(), "breaking")
	breakingIndicator := parsedMessageBody[breakingIndicatorIndex]
	if breakingIndicator == "!" {
		return BreakingChange, nil
	}

	typeIndex := util.MustFind(bodyRegex.SubexpNames(), "type")

	return classifier.StringToType(parsedMessageBody[typeIndex])
}
