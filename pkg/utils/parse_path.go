package utils

import (
	"fmt"
	"strings"
)

func ParsePathWithUName(path string, username string) (string, []string, error) {
	if username == "" {
		return "", nil, fmt.Errorf("username cannot be empty")
	}

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	path_slice := strings.Split(path, "/")

	entryName := path_slice[len(path_slice)-1]
	path_slice = path_slice[:len(path_slice)-1]

	path_slice = append([]string{username}, path_slice...)

	return entryName, path_slice, nil
}
