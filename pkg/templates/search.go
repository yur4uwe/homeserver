package templates

type SearchEntry struct {
	Name     string
	Path     string
	IsFolder bool
	Size     string
}

type SearchPageData struct {
	Query     string
	Directory string
	Parent    string
	Entries   []SearchEntry
}
