package templates

type UploadPageData struct {
	Files     []FileInfo
	FileCount int
}

type FileInfo struct {
	Name        string
	Size        int64
	ContentType string
}
