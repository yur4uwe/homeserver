package globals

var Pages []string
var ProtectedPages []string

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
		if page == path {
			return true
		}
	}
	return false
}

func IsProtectedPath(path string) bool {
	for _, page := range ProtectedPages {
		if page == path {
			return true
		}
	}
	return false
}
