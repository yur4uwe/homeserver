package file_schema

import (
	"fmt"
	"strings"
)

func CheckFileID(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.ID == 0 {
		return uf, fmt.Errorf("file ID is invalid: ID must be greater than 0")
	}
	return uf, nil
}

func CheckFilename(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.Filename == "" {
		return uf, fmt.Errorf("filename is invalid: filename cannot be empty")
	}
	return uf, nil
}

func CheckExtension(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.Extension == "" {
		return uf, fmt.Errorf("extension is invalid: extension cannot be empty")
	}
	if strings.HasPrefix(uf.Extension, ".") {
		return uf, fmt.Errorf("extension is invalid: extension cannot start with a dot")
	}
	return uf, nil
}

func CheckPath(ownerName string) func(uf *UserFile, err error) (*UserFile, error) {
	return func(uf *UserFile, err error) (*UserFile, error) {
		if err != nil {
			return uf, err
		}
		if uf.Path == "" {
			return uf, fmt.Errorf("path is invalid: path cannot be empty")
		}
		if strings.HasPrefix(uf.Path, "/") {
			return uf, fmt.Errorf("path is invalid: path cannot start with a slash")
		}
		if strings.HasSuffix(uf.Path, "/") {
			return uf, fmt.Errorf("path is invalid: path cannot end with a slash")
		}
		if !strings.HasPrefix(uf.Path, ownerName) {
			return uf, fmt.Errorf("path is invalid: path must start with the owner's name '%s'", ownerName)
		}
		return uf, nil
	}

}

func CheckDirectory(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.Directory == 0 {
		return uf, fmt.Errorf("directory is invalid: directory ID must be greater than 0")
	}
	return uf, nil
}

func CheckSize(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.Size <= 0 {
		return uf, fmt.Errorf("size is invalid: file size must be greater than 0")
	}
	return uf, nil
}

func CheckContentType(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.ContentType == "" {
		return uf, fmt.Errorf("content type is invalid: content type cannot be empty")
	}
	return uf, nil
}

func CheckOwnerID(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.OwnerID == 0 {
		return uf, fmt.Errorf("owner ID is invalid: owner ID must be greater than 0")
	}
	return uf, nil
}

func CheckMetadata(uf *UserFile, err error) (*UserFile, error) {
	if err != nil {
		return uf, err
	}
	if uf.Metadata == nil {
		return uf, fmt.Errorf("metadata is invalid: metadata cannot be nil")
	}
	if *uf.Metadata == "" {
		return uf, fmt.Errorf("metadata is invalid: metadata cannot be empty")
	}
	return uf, nil
}
