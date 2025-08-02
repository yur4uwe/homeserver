package globals

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	EnvPort          = "PORT"
	EnvDbURL         = "DB_URL"
	EnvStoragePath   = "STORAGE_PATH"
	EnvEmailPassword = "EMAIL_PASSWORD"
	EnvEmailHost     = "EMAIL_HOST"
	EnvEmailPort     = "EMAIL_PORT"
	EnvEmailUser     = "EMAIL_USER"
)

var Pages []string
var ProtectedPages []string
var AvailableFileTypes = []string{}

func FillPages() {
	Pages = []string{
		"/login",
		"/register",
		"/verification",
		"/settings",
		"/upload",
		"/search/",
	}

	ProtectedPages = []string{
		"/settings",
		"/upload",
		"/search/",
	}

	data, err := os.ReadFile("../resources/catalog.json")
	if err != nil {
		panic("unable to read catalog.json: " + err.Error())
	}

	err = json.Unmarshal(data, &AvailableFileTypes)
	if err != nil {
		panic("unable to parse catalog.json: " + err.Error())
	}
}

func IsPage(path string) bool {
	if path == "/" {
		return true
	}

	for _, page := range Pages {
		if page[len(page)-1] == '/' && len(page) > 1 && strings.HasPrefix(path, page) {
			fmt.Println("Matched path:", path, "with prefix:", page)
			return true
		}

		if path == page {
			fmt.Println("Matched exact path:", path, "with page:", page)
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

func BuildFilePath(relative_path, filename, extension string) string {
	storage_path := os.Getenv(EnvStoragePath)
	if storage_path[len(storage_path)-1] != '/' {
		panic("Storage path must end with a slash")
	}

	if len(relative_path) > 0 && (relative_path[len(relative_path)-1] == '/' || relative_path[0] == '/') {
		panic("Relative path must not start or end with a slash")
	}

	var file_path string
	if relative_path != "" {
		file_path = fmt.Sprintf("%s/%s", relative_path, filename)
	} else {
		file_path = filename
	}

	return fmt.Sprintf("%s%s.%s", storage_path, file_path, extension)
}
