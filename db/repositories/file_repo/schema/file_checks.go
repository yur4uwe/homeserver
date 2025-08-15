package fschema

import (
	"fmt"
	"homeserver/constants"
	"homeserver/constants/fconsts"
	"slices"
	"strings"
)

func (uf *UserFile) ShouldHaveFields(ownerName string) func(...string) error {
	return func(fields ...string) error {
		if len(fields) == 0 {
			return fmt.Errorf(fconsts.NoFieldsSpecifiedError)
		}

		var errorMessages []string

		for _, field := range fields {
			if !slices.Contains(constants.DBFilesTable.Fields, field) {
				errorMessages = append(errorMessages, fmt.Sprintf(fconsts.InvalidFileFieldError, field))
				continue
			}

			switch field {
			case fconsts.ID:
				if err := CheckFileID(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("FileID: %v", err))
				}
			case fconsts.Filename:
				if err := CheckFilename(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Filename: %v", err))
				}
			case fconsts.Extension:
				if err := CheckExtension(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Extension: %v", err))
				}
			case fconsts.Path:
				if err := CheckPath(ownerName)(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Path: %v", err))
				}
			case fconsts.Directory:
				if err := CheckDirectory(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Directory: %v", err))
				}
			case fconsts.Size:
				if err := CheckSize(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Size: %v", err))
				}
			case fconsts.ContentType:
				if err := CheckContentType(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("ContentType: %v", err))
				}
			case fconsts.OwnerID:
				if err := CheckOwnerID(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("OwnerID: %v", err))
				}
			case fconsts.Metadata:
				if err := CheckMetadata(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Metadata: %v", err))
				}
			}
		}

		if len(errorMessages) > 0 {
			return fmt.Errorf("file validation failed:\n%s", strings.Join(errorMessages, "\n"))
		}

		return nil
	}
}

func CheckFileID(uf *UserFile) error {
	if uf.ID == 0 {
		return fmt.Errorf(fconsts.InvalidIDError)
	}
	return nil
}

func CheckFilename(uf *UserFile) error {
	if uf.Filename == "" {
		return fmt.Errorf(fconsts.MissingNameError)
	}
	return nil
}

func CheckExtension(uf *UserFile) error {
	if uf.Extension == "" {
		return fmt.Errorf(fconsts.InvalidExtensionError)
	}
	if strings.HasPrefix(uf.Extension, ".") {
		return fmt.Errorf(fconsts.InvalidExtensionError)
	}
	return nil
}

func CheckPath(ownerName string) func(uf *UserFile) error {
	return func(uf *UserFile) error {
		if uf.Path == "" {
			return fmt.Errorf(fconsts.MissingPathError)
		}
		if strings.HasPrefix(uf.Path, "/") {
			return fmt.Errorf(fconsts.InvalidPathStartError)
		}
		if strings.HasSuffix(uf.Path, "/") {
			return fmt.Errorf(fconsts.InvalidPathEndError)
		}
		if !strings.HasPrefix(uf.Path, ownerName) {
			return fmt.Errorf(fconsts.MissingOwnerFolderError)
		}
		return nil
	}
}

func CheckDirectory(uf *UserFile) error {
	if uf.Directory == 0 {
		return fmt.Errorf(fconsts.InvalidDirectoryError)
	}
	return nil
}

func CheckSize(uf *UserFile) error {
	if uf.Size <= 0 {
		return fmt.Errorf(fconsts.InvalidSizeError)
	}
	return nil
}

func CheckContentType(uf *UserFile) error {
	if uf.ContentType == "" {
		return fmt.Errorf(fconsts.InvalidContentTypeError)
	}
	return nil
}

func CheckOwnerID(uf *UserFile) error {
	if uf.OwnerID == 0 {
		return fmt.Errorf(fconsts.InvalidOwnerIDError)
	}
	return nil
}

func CheckMetadata(uf *UserFile) error {
	if uf.Metadata == nil || *uf.Metadata == "" {
		return fmt.Errorf(fconsts.MissingMetadataError)
	}
	return nil
}
