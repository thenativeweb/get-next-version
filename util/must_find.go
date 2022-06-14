package util

func MustFind[T comparable](slice []T, itemToFind T) int {
	for index, item := range slice {
		if item == itemToFind {
			return index
		}
	}

	panic("item not found")
}
