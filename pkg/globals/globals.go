package globals

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
