package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/get-next-version/util"
)

func TestMustFind(t *testing.T) {
	tests := []struct {
		slice         []string
		itemToFind    string
		willPanic     bool
		expectedIndex int
	}{
		{slice: []string{}, itemToFind: "some-item", willPanic: true, expectedIndex: 0},
		{slice: []string{"some-item"}, itemToFind: "some-item", willPanic: false, expectedIndex: 0},
		{slice: []string{"other-item", "some-item"}, itemToFind: "some-item", willPanic: false, expectedIndex: 1},
		{slice: []string{"other-item", "some-item", "some-item"}, itemToFind: "some-item", willPanic: false, expectedIndex: 1},
	}

	for _, test := range tests {
		if test.willPanic {
			assert.Panics(t, func() {
				util.MustFind(test.slice, test.itemToFind)
			})
			continue
		}

		index := util.MustFind(test.slice, test.itemToFind)
		assert.Equal(t, test.expectedIndex, index)
	}
}
