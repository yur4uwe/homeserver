package fschema

import (
	"os"
	"strings"
	"time"
)

var (
	DefaultFile = UserFile{
		ID:            0,
		Filename:      "",
		Extension:     "",
		Path:          "",
		Directory:     0,
		Size:          0,
		ContentType:   "",
		OwnerID:       0,
		TrustedSource: []byte{},
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		Metadata:      nil,
	}
)

var (
	Getenv = os.Getenv
)

func NewUserFileBuilder() *UserFile {
	clone := DefaultFile
	return &clone
}

func (uf *UserFile) WithID(id int) *UserFile {
	if id <= 0 {
		panic("File ID must be greater than 0")
	}
	uf.ID = id
	return uf
}

func (uf *UserFile) WithFilename(filename string) *UserFile {
	if filename == "" {
		panic("Filename cannot be empty")
	}
	uf.Filename = filename
	return uf
}

func (uf *UserFile) WithFullName(fullName string) *UserFile {
	if fullName == "" {
		panic("Full name cannot be empty")
	}
	parts := strings.Split(fullName, ".")
	uf.Filename = strings.Join(parts[:len(parts)-1], ".")
	uf.Extension = parts[len(parts)-1]
	return uf
}

func (uf *UserFile) WithExtension(extension string) *UserFile {
	if extension == "" {
		panic("Extension cannot be empty")
	}
	if strings.HasPrefix(extension, ".") {
		panic("Extension cannot start with a dot")
	}
	uf.Extension = extension
	return uf
}

func (uf *UserFile) WithPath(path string) *UserFile {
	if path == "" {
		panic("Path cannot be empty")
	}
	if path[0] == '/' || path[len(path)-1] == '/' {
		panic("Path cannot start or end with a slash")
	}
	uf.Path = path
	return uf
}

func (uf *UserFile) WithDirectory(directory int) *UserFile {
	if directory <= 0 {
		panic("Directory ID must be greater than 0")
	}
	uf.Directory = directory
	return uf
}

func (uf *UserFile) WithSize(size int64) *UserFile {
	if size <= 0 {
		panic("File size must be greater than 0")
	}
	uf.Size = size
	return uf
}

func (uf *UserFile) WithContentType(contentType string) *UserFile {
	if contentType == "" {
		panic("Content type cannot be empty")
	}
	uf.ContentType = contentType
	return uf
}

func (uf *UserFile) WithOwnerID(ownerID int) *UserFile {
	if ownerID <= 0 {
		panic("Owner ID must be greater than 0")
	}
	uf.OwnerID = ownerID
	return uf
}

func (uf *UserFile) WithMetadata(metadata *string) *UserFile {
	if metadata == nil || *metadata == "" {
		panic("Metadata cannot be nil or empty")
	}
	uf.Metadata = metadata
	return uf
}

func (uf *UserFile) WithTrustedSource(trustedSource []byte) *UserFile {
	uf.TrustedSource = trustedSource
	return uf
}
