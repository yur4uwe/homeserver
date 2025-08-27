package templates

type SearchEntry struct {
	ID        int
	Name      string
	URIName   string // e.g., "example_file"
	Extension string // e.g., ".txt", ".jpg"
	// URIName is the name of the file or folder as it appears in the URL,
	Path     string // e.g., "files/example_file.txt"
	URIPath  string // e.g., "files/example_file.txt"
	Type     string // e.g., "svg", "txt", "jpg"
	IsFolder bool
	Size     string // e.g., "1.2 MB"
}

type SearchPageData struct {
	Query        string
	URIQuery     string
	Directory    string
	URIDirectory string
	Parent       string
	URIParent    string
	DownloadPath string
	QueryResults bool
	Entries      []SearchEntry
}
