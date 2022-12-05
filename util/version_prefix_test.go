package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/util"
)

func TestIsValidVersionPrefix(t *testing.T) {
	for _, testcase := range []struct {
		name     string
		prefix   string
		expected bool
	}{
		{name: "empty prefix", prefix: "", expected: true},
		{name: "valid v prefix", prefix: "v", expected: true},
		{name: "trailing digit", prefix: "test1", expected: false},

		{name: "prefix containing double slashes", prefix: "test/foo/bar", expected: true},
		{name: "prefix containing double consecutive slashes", prefix: "test//foo", expected: false},

		{name: "prefix ending with dot", prefix: "test.", expected: true},
		{name: "prefix containing a single dot", prefix: "test..foo", expected: false},
		{name: "prefix containing duplicated dots", prefix: "test..foo", expected: false},

		{name: "tag name starting with slash", prefix: "/test", expected: false},
		{name: "tag name ending with slash", prefix: "test/", expected: true},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			result, err := util.IsValidVersionPrefix(testcase.prefix)
			if !testcase.expected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testcase.expected, result)
		})
	}
}
