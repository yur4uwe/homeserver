package globals

var Pages []string
var ProtectedPages []string

func FillPages() {
	Pages = []string{
		"/",
		"/login",
		"/register",
	}

	ProtectedPages = []string{}
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
