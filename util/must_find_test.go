package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/getnextversion/util"
)

func TestMustFind(t *testing.T) {
	tests := []struct {
		slice         []string
		itemToFind    string
		withError     bool
		expectedIndex int
	}{
		{slice: []string{}, itemToFind: "some-item", withError: true, expectedIndex: 0},
		{slice: []string{"some-item"}, itemToFind: "some-item", withError: false, expectedIndex: 0},
		{slice: []string{"other-item", "some-item"}, itemToFind: "some-item", withError: false, expectedIndex: 1},
		{slice: []string{"other-item", "some-item", "some-item"}, itemToFind: "some-item", withError: false, expectedIndex: 1},
	}

	for _, test := range tests {
		index, err := util.MustFind(test.slice, test.itemToFind)

		if test.withError {
			assert.Error(t, err)
			continue
		}
		assert.Equal(t, test.expectedIndex, index)
	}
}
