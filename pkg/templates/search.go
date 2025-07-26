package templates

type SearchEntry struct {
	Name     string
	Path     string
	IsFolder bool
}

type SearchPageData struct {
	Query     string
	Directory string
	Entries   []SearchEntry
}
