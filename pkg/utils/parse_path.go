package utils

import (
	"fmt"
	"strings"
)

func ParsePathWithUName(path string, username string) ([]string, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	path_slice := strings.Split(path, "/")

	path_slice = append([]string{username}, path_slice...)

	return path_slice, nil
}
