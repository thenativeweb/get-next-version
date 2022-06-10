package util

import "strings"

type IsOnePrefixResult struct {
	IsOnePrefix bool
	Prefix      string
}

func IsOnePrefix(s string, possiblePrefixes []string) IsOnePrefixResult {
	for _, possiblePrefix := range possiblePrefixes {
		if strings.HasPrefix(s, possiblePrefix) {
			return IsOnePrefixResult{
				IsOnePrefix: true,
				Prefix:      possiblePrefix,
			}
		}
	}

	return IsOnePrefixResult{
		IsOnePrefix: false,
		Prefix:      "",
	}
}
