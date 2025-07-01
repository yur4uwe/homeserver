package globals

import (
	"fmt"
	"os"
	"strings"
)

var Pages []string

func FillPages() {
	Pages = []string{
		"/",
		"/login",
		"/register",
	}
}

func IsPage(path string) bool {
	for _, page := range Pages {
		if page == path {
			return true
		}
	}
	return false
}

func MergePage(newPageContent string) ([]byte, error) {
	file, err := os.ReadFile("../static/index.html")
	if err != nil {
		return nil, err
	}

	fileContent := string(file)
	if newPageContent == "" {
		return []byte(fileContent), nil
	}

	slotIndex := strings.Index(fileContent, "<slot></slot>") + len("<slot>")
	if slotIndex == len("<slot>")-1 {
		return nil, fmt.Errorf("slot not found in index.html")
	}

	mergedContent := fileContent[:slotIndex] + "\n" + newPageContent + "\n" + fileContent[slotIndex:]

	return []byte(mergedContent), nil
}
