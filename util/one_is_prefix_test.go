package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/util"
)

func TestOneIsPrefix(t *testing.T) {
	tests := []struct {
		string              string
		possiblePrefixes    []string
		expectedIsOnePrefix bool
		expectedPrefix      string
	}{
		{string: "foo", possiblePrefixes: []string{}, expectedIsOnePrefix: false, expectedPrefix: ""},
		{string: "foo", possiblePrefixes: []string{"o"}, expectedIsOnePrefix: false, expectedPrefix: ""},
		{string: "foo", possiblePrefixes: []string{"f"}, expectedIsOnePrefix: true, expectedPrefix: "f"},
		{string: "foo", possiblePrefixes: []string{"o", "f"}, expectedIsOnePrefix: true, expectedPrefix: "f"},
		{string: "foo", possiblePrefixes: []string{"o", "f", "fo"}, expectedIsOnePrefix: true, expectedPrefix: "f"},
	}

	for _, test := range tests {
		actual := util.IsOnePrefix(test.string, test.possiblePrefixes)

		assert.Equal(t, test.expectedIsOnePrefix, actual.IsOnePrefix)
		assert.Equal(t, test.expectedPrefix, actual.Prefix)
	}
}
