package constants

import (
	"encoding/json"
	"os"
)

var AvailableFileTypes = []string{}

func ReadAvailableFileTypes() {
	data, err := os.ReadFile("../resources/catalog.json")
	if err != nil {
		panic("unable to read catalog.json: " + err.Error())
	}

	err = json.Unmarshal(data, &AvailableFileTypes)
	if err != nil {
		panic("unable to parse catalog.json: " + err.Error())
	}
}
