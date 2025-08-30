package constants

var (
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
		"/api/v1/settings",
		"/upload",
		"/search/",
	}
)
