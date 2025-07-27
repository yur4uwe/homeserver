package templates

type SearchEntry struct {
	Name    string
	URIName string // e.g., "example_file.txt"
	// URIName is the name of the file or folder as it appears in the URL,
	Path     string
	URIPath  string // e.g., "/files/example_file.txt"
	Type     string // e.g., "svg", "txt", "jpg"
	IsFolder bool
	Size     string
}

type SearchPageData struct {
	Query        string
	URIQuery     string
	Directory    string
	URIDirectory string
	Parent       string
	URIParent    string
	Entries      []SearchEntry
}
