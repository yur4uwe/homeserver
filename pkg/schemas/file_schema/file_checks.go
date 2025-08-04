package file_schema

import (
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas"
	"slices"
	"strings"
)

func (uf *UserFile) ShouldHaveFields(ownerName string) func(...string) error {
	return func(fields ...string) error {
		if len(fields) == 0 {
			return fmt.Errorf(schemas.NoFieldsSpecifiedError)
		}

		var errorMessages []string

		for _, field := range fields {
			if !slices.Contains(globals.DBFilesTable.Fields, field) {
				errorMessages = append(errorMessages, fmt.Sprintf(schemas.InvalidFileFieldError, field))
				continue
			}

			switch field {
			case schemas.FileID:
				if err := CheckFileID(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("FileID: %v", err))
				}
			case schemas.Filename:
				if err := CheckFilename(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Filename: %v", err))
				}
			case schemas.Extension:
				if err := CheckExtension(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Extension: %v", err))
				}
			case schemas.Path:
				if err := CheckPath(ownerName)(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Path: %v", err))
				}
			case schemas.Directory:
				if err := CheckDirectory(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Directory: %v", err))
				}
			case schemas.Size:
				if err := CheckSize(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Size: %v", err))
				}
			case schemas.ContentType:
				if err := CheckContentType(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("ContentType: %v", err))
				}
			case schemas.OwnerID:
				if err := CheckOwnerID(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("OwnerID: %v", err))
				}
			case schemas.Metadata:
				if err := CheckMetadata(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Metadata: %v", err))
				}
			}
		}

		if len(errorMessages) > 0 {
			return fmt.Errorf("%s%s", schemas.ValidationErrorPrefix, strings.Join(errorMessages, "\n"))
		}

		return nil
	}
}

func CheckFileID(uf *UserFile) error {
	if uf.ID == 0 {
		return fmt.Errorf(schemas.InvalidFileIDError)
	}
	return nil
}

func CheckFilename(uf *UserFile) error {
	if uf.Filename == "" {
		return fmt.Errorf(schemas.InvalidFilenameError)
	}
	return nil
}

func CheckExtension(uf *UserFile) error {
	if uf.Extension == "" {
		return fmt.Errorf(schemas.InvalidExtensionError)
	}
	if strings.HasPrefix(uf.Extension, ".") {
		return fmt.Errorf(schemas.InvalidExtensionError)
	}
	return nil
}

func CheckPath(ownerName string) func(uf *UserFile) error {
	return func(uf *UserFile) error {
		if uf.Path == "" {
			return fmt.Errorf(schemas.InvalidFolderPathError)
		}
		if strings.HasPrefix(uf.Path, "/") {
			return fmt.Errorf(schemas.InvalidFilePathError)
		}
		if strings.HasSuffix(uf.Path, "/") {
			return fmt.Errorf(schemas.InvalidFilePathEndError)
		}
		if !strings.HasPrefix(uf.Path, ownerName) {
			return fmt.Errorf(schemas.InvalidFilePathOwnerError, ownerName)
		}
		return nil
	}
}

func CheckDirectory(uf *UserFile) error {
	if uf.Directory == 0 {
		return fmt.Errorf(schemas.InvalidDirectoryError)
	}
	return nil
}

func CheckSize(uf *UserFile) error {
	if uf.Size <= 0 {
		return fmt.Errorf(schemas.InvalidFileSizeError)
	}
	return nil
}

func CheckContentType(uf *UserFile) error {
	if uf.ContentType == "" {
		return fmt.Errorf(schemas.InvalidContentTypeError)
	}
	return nil
}

func CheckOwnerID(uf *UserFile) error {
	if uf.OwnerID == 0 {
		return fmt.Errorf(schemas.InvalidFileOwnerIDError)
	}
	return nil
}

func CheckMetadata(uf *UserFile) error {
	if uf.Metadata == nil || *uf.Metadata == "" {
		return fmt.Errorf(schemas.InvalidFileMetadataError)
	}
	return nil
}
