package templates

type SearchEntry struct {
	Name     string
	Path     string
	Type     string // e.g., ".svg", ".txt", ".jpg"
	IsFolder bool
	Size     string
}

type SearchPageData struct {
	Query     string
	Directory string
	Parent    string
	Entries   []SearchEntry
}
