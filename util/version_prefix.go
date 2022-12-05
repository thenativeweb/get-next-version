package util

import (
	"errors"
	"regexp"
)

var (
	versionPrefixCharSetValidation         *regexp.Regexp
	versionPrefixConsecutiveCharValidation *regexp.Regexp
	versionPrefixLastCharValidation        *regexp.Regexp
)

func init() {
	versionPrefixCharSetValidation = regexp.MustCompile(`^[^/][a-zA-Z\d\.\-_/]*$`)
	versionPrefixConsecutiveCharValidation = regexp.MustCompile(`/{2,}|\.{2,}`)
	versionPrefixLastCharValidation = regexp.MustCompile(`\D$`)
}

/*
IsValidVersionPrefix runs a series of checks to ensure that the version prefix is valid
Allowed characters: alphanumerics, dot, dash, slash, underscore
Version prefixes must not:
  - begin with a slash
  - include multiple consecutive slashes and dots
  - end with a digit
*/
func IsValidVersionPrefix(prefix string) (bool, error) {
	if prefix == "" {
		return true, nil
	}

	if !versionPrefixCharSetValidation.MatchString(prefix) {
		return false, errors.New("version prefix must not contain invalid character(s) or begin with slash")
	}
	if versionPrefixConsecutiveCharValidation.MatchString(prefix) {
		return false, errors.New("version prefix must not contain multiple consecutive slashes or dots")
	}
	if !versionPrefixLastCharValidation.MatchString(prefix) {
		return false, errors.New("version prefix must not end with a digit")
	}

	return true, nil
}
