package globals

import "strings"

var Pages []string
var ProtectedPages []string
var AvailableFileTypes = []string{
	".svg",
	".txt",
	".jpg",
	".png",
	".pdf",
	".docx",
	".xlsx",
	".pptx",
	".mp3",
	".mp4",
	".avi",
	".zip",
	".rar",
}

func FillPages() {
	Pages = []string{
		"/",
		"/login",
		"/register",
		"/verification",
		"/settings",
		"/upload",
		"/search",
	}

	ProtectedPages = []string{
		"/settings",
		"/upload",
		"/search",
	}
}

func IsPage(path string) bool {
	for _, page := range Pages {
		if strings.HasPrefix(path, page) {
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
