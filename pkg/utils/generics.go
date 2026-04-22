package utils

import "slices"

func Without[T comparable](slice []T, items ...T) []T {
	var result []T
	for _, item := range slice {
		if !slices.Contains(items, item) {
			result = append(result, item)
		}
	}
	return result
}

func Must(value interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return value
}
