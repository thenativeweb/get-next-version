package util

import "errors"

func MustFind[T comparable](slice []T, itemToFind T) (int, error) {
	for index, item := range slice {
		if item == itemToFind {
			return index, nil
		}
	}

	return 0, errors.New("item not found")
}
