package file_schema

import (
	"homeserver/pkg/globals"
	"time"
)

func NewUserFileBuilder() *UserFile {
	return &UserFile{
		ID:          0,
		Filename:    "",
		Extension:   "",
		Path:        "",
		Directory:   0,
		Size:        0,
		ContentType: "",
		OwnerID:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    nil,
	}
}

func (uf *UserFile) WithID(id int) *UserFile {
	uf.ID = id
	return uf
}

func (uf *UserFile) WithFilename(filename string) *UserFile {
	uf.Filename = filename
	return uf
}

func (uf *UserFile) WithExtension(extension string) *UserFile {
	uf.Extension = extension
	return uf
}

func (uf *UserFile) WithPath(path string) *UserFile {
	uf.Path = path
	return uf
}

func (uf *UserFile) WithDirectory(directory int) *UserFile {
	uf.Directory = directory
	return uf
}

func (uf *UserFile) WithSize(size int64) *UserFile {
	uf.Size = size
	return uf
}

func (uf *UserFile) WithContentType(contentType string) *UserFile {
	uf.ContentType = contentType
	return uf
}

func (uf *UserFile) WithOwnerID(ownerID int) *UserFile {
	uf.OwnerID = ownerID
	return uf
}

func (uf *UserFile) WithMetadata(metadata *string) *UserFile {
	uf.Metadata = metadata
	return uf
}

func (uf *UserFile) AbsolutePath() string {
	return globals.BuildFilePath(uf.Path, uf.Filename, uf.Extension)
}
