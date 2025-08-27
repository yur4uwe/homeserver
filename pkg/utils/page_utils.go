package utils

import (
	"homeserver/constants"
	"strings"
)

func IsPage(path string) bool {
	if path == "/" {
		return true
	}

	for _, page := range constants.Pages {
		if page[len(page)-1] == '/' && len(page) > 1 && strings.HasPrefix(path, page) {
			return true
		}

		if path == page {
			return true
		}
	}

	return false
}

func IsProtectedPath(path string) bool {
	for _, page := range constants.ProtectedPages {
		if strings.HasPrefix(path, page) {
			return true
		}
	}
	return false
}
