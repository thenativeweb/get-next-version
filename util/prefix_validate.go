package util

import (
	"errors"
	"regexp"
)

/*
IsValidTagPrefix runs a series of regex checks to ensure that the tag prefix is valid
  Tags cannot begin or end with, or contain multiple consecutive / characters.
  They cannot contain any of the following characters \, ?, ~, ^, :, * , [, @.
  They cannot contain a space.
  They cannot end with a . or have two consecutive .. anywhere within them.
  Tags are not case-sensitive.
*/
func IsValidTagPrefix(prefix string) (bool, error) {
	if prefix == "" {
		return true, nil
	}

	if !regexp.MustCompile(`^[^/][a-zA-Z\d\.\-_/]*$`).MatchString(prefix) {
		return false, errors.New("tag prefix contains invalid character(s) or begins with /")
	}
	if regexp.MustCompile(`/{2,}`).MatchString(prefix) {
		return false, errors.New("tag prefix contains multiple consecutive / characters")
	}
	if regexp.MustCompile(`\.{2,}`).MatchString(prefix) {
		return false, errors.New("tag prefix contains multiple consecutive . characters")
	}
	if !regexp.MustCompile(`\D$`).MatchString(prefix) {
		return false, errors.New("tag prefix ends with a digit")
	}

	return true, nil
}
