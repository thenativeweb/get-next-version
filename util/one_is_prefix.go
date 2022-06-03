package util

import "strings"

func OneIsPrefix(s string, possiblePrefixes []string) bool {
	for _, possiblePrefix := range possiblePrefixes {
		if strings.HasPrefix(s, possiblePrefix) {
			return true
		}
	}

	return false
}
