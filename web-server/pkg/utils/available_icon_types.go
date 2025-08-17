package utils

import (
	"homeserver/constants"
	"strings"
)

func IsFileTypeAvailable(fileType string) bool {
	for _, availableType := range constants.AvailableFileTypes {
		if strings.EqualFold(availableType, fileType) {
			return true
		}
	}
	return false
}
