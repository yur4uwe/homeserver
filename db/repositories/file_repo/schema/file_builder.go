package fschema

import (
	"fmt"
	"homeserver/constants"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"homeserver/web-server/pkg/utils"
	"os"
	"strings"
	"time"
)

var (
	defaultFile = UserFile{
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
)

var (
	Getenv = os.Getenv
)

func NewUserFileBuilder() *UserFile {
	clone := defaultFile
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

func (uf *UserFile) Get(name string) (any, error) {
	switch name {
	case fconsts.ID:
		return uf.ID, nil
	case fconsts.Filename:
		return uf.Filename, nil
	case fconsts.Extension:
		return uf.Extension, nil
	case fconsts.Path:
		return uf.Path, nil
	case fconsts.Directory:
		return uf.Directory, nil
	case fconsts.Size:
		return uf.Size, nil
	case fconsts.ContentType:
		return uf.ContentType, nil
	case fconsts.OwnerID:
		return uf.OwnerID, nil
	case fconsts.CreatedAt:
		return uf.CreatedAt, nil
	case fconsts.UpdatedAt:
		return uf.UpdatedAt, nil
	default:
		return nil, fmt.Errorf("unknown field: %s", name)
	}
}

func (uf *UserFile) Set(name string, value any) error {
	switch name {
	case fconsts.ID:
		if id, ok := value.(int); ok && id > 0 {
			uf.ID = id
			return nil
		}
		return fmt.Errorf("invalid ID value: %v", value)
	case fconsts.Filename:
		if filename, ok := value.(string); ok && filename != "" {
			uf.Filename = filename
			return nil
		}
		return fmt.Errorf("invalid Filename value: %v", value)
	case fconsts.Extension:
		if ext, ok := value.(string); ok && ext != "" {
			uf.Extension = ext
			return nil
		}
		return fmt.Errorf("invalid Extension value: %v", value)
	case fconsts.Path:
		if path, ok := value.(string); ok && path != "" {
			uf.Path = path
			return nil
		}
		return fmt.Errorf("invalid Path value: %v", value)
	case fconsts.Directory:
		if dir, ok := value.(int); ok && dir > 0 {
			uf.Directory = dir
			return nil
		}
		return fmt.Errorf("invalid Directory value: %v", value)
	case fconsts.Size:
		if size, ok := value.(int64); ok && size > 0 {
			uf.Size = size
			return nil
		}
		return fmt.Errorf("invalid Size value: %v", value)
	case fconsts.ContentType:
		if contentType, ok := value.(string); ok && contentType != "" {
			uf.ContentType = contentType
			return nil
		}
		return fmt.Errorf("invalid ContentType value: %v", value)
	case fconsts.OwnerID:
		if ownerID, ok := value.(int); ok && ownerID > 0 {
			uf.OwnerID = ownerID
			return nil
		}
		return fmt.Errorf("invalid OwnerID value: %v", value)
	default:
		return fmt.Errorf("unknown field: %s", name)
	}
}

func (uf *UserFile) AbsolutePath() string {
	storage_path := Getenv(constants.StoragePath)

	var file_path string
	if uf.Path != "" {
		file_path = fmt.Sprintf("%s/%s", uf.Path, uf.Filename)
	} else {
		file_path = uf.Filename
	}

	return fmt.Sprintf("%s%s.%s", storage_path, file_path, uf.Extension)
}

func FieldsList(without ...string) []string {
	return utils.Without(dbconsts.DBFilesTable.Fields, without...)
}
