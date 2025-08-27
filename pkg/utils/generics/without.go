package generics

import "slices"

func Without[T comparable](slice []T, items ...T) []T {
	var result []T
	for _, item := range slice {
		if !contains(item, items) {
			result = append(result, item)
		}
	}
	return result
}

func contains[T comparable](item T, slice []T) bool {
	return slices.Contains(slice, item)
}
