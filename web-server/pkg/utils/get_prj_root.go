package utils

import (
	"os"
	"strings"
)

const (
	root_fldr_name = "homeserver"
)

func GetProjectRoot() string {
	wd, _ := os.Getwd()

	strs := strings.Split(wd, string(os.PathSeparator))

	for i := len(strs) - 1; i >= 0; i-- {
		if strs[i] == root_fldr_name {
			return strings.Join(strs[:i+1], string(os.PathSeparator)) + string(os.PathSeparator)
		}
	}

	return ""
}
