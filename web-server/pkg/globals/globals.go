package globals

import (
	"strings"
)

func IsPage(path string) bool {
	if path == "/" {
		return true
	}

	for _, page := range Pages {
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
	for _, page := range ProtectedPages {
		if strings.HasPrefix(path, page) {
			return true
		}
	}
	return false
}

func IsFileTypeAvailable(fileType string) bool {
	for _, availableType := range AvailableFileTypes {
		if strings.EqualFold(availableType, fileType) {
			return true
		}
	}
	return false
}
