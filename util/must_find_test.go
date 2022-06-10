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
		doPanic       bool
		expectedIndex int
	}{
		{slice: []string{}, itemToFind: "some-item", doPanic: true, expectedIndex: 0},
		{slice: []string{"some-item"}, itemToFind: "some-item", doPanic: false, expectedIndex: 0},
		{slice: []string{"other-item", "some-item"}, itemToFind: "some-item", doPanic: false, expectedIndex: 1},
		{slice: []string{"other-item", "some-item", "some-item"}, itemToFind: "some-item", doPanic: false, expectedIndex: 1},
	}

	for _, test := range tests {
		if test.doPanic {
			assert.Panics(t, func() { util.MustFind(test.slice, test.itemToFind) })
			continue
		}

		assert.Equal(t, test.expectedIndex, util.MustFind(test.slice, test.itemToFind))
	}
}
