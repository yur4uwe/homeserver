package utils

import (
	"encoding/json"
	"fmt"
	"homeserver/constants"
	"os"
	"strings"
)

func LoadAvailableIconTypes() error {
	f, err := os.ReadFile(constants.Rootd + "resources/catalog.json")
	if err != nil {
		return fmt.Errorf("error reading catalog.json: %v", err)
	}

	var iconTypes []string
	if err := json.Unmarshal(f, &iconTypes); err != nil {
		return fmt.Errorf("error unmarshalling catalog.json: %v", err)
	}

	constants.AvailableFileTypes = iconTypes

	return nil
}

func IsFileTypeAvailable(fileType string) bool {
	for _, availableType := range constants.AvailableFileTypes {
		if strings.EqualFold(availableType, fileType) {
			return true
		}
	}
	return false
}
