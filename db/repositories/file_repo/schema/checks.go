package fschema

import (
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"slices"
	"strings"
)

func (uf *UserFile) checkFields(ownerName string, fields []string, shouldHave bool) error {
	if len(fields) == 0 {
		return fmt.Errorf(fconsts.NoFieldsSpecifiedError)
	}

	var errorMessages []string

	for _, field := range fields {
		if !slices.Contains(dbconsts.DBFilesTable.Fields, field) {
			errorMessages = append(errorMessages, fmt.Sprintf(fconsts.InvalidFileFieldError, field))
			continue
		}

		var err error
		switch field {
		case fconsts.ID:
			err = CheckFileID(uf)
		case fconsts.Filename:
			err = CheckFilename(uf)
		case fconsts.Extension:
			err = CheckExtension(uf)
		case fconsts.Path:
			err = CheckPath(ownerName)(uf)
		case fconsts.Directory:
			err = CheckDirectory(uf)
		case fconsts.Size:
			err = CheckSize(uf)
		case fconsts.ContentType:
			err = CheckContentType(uf)
		case fconsts.OwnerID:
			err = CheckOwnerID(uf)
		case fconsts.Metadata:
			err = CheckMetadata(uf)
		case fconsts.TrustedSource:
			err = CheckTrustedSource(uf)
		case fconsts.CreatedAt:
			if uf.CreatedAt.IsZero() {
				err = fmt.Errorf(fconsts.InvalidCreatedAtError)
			}
		case fconsts.UpdatedAt:
			if uf.UpdatedAt.IsZero() {
				err = fmt.Errorf(fconsts.InvalidUpdatedAtError)
			}
		default:
			err = fmt.Errorf(fconsts.InvalidFileFieldError, field)
		}

		if shouldHave {
			if err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("%s: %v", field, err))
			}
		} else {
			if err == nil {
				errorMessages = append(errorMessages, fmt.Sprintf("%s should not be set", field))
			}
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("file validation failed:\n%s", strings.Join(errorMessages, "\n"))
	}
	return nil
}

func (uf *UserFile) ShouldHaveFields(ownerName string) func(...string) error {
	return func(fields ...string) error {
		return uf.checkFields(ownerName, fields, true)
	}
}

func (uf *UserFile) ShouldNotHaveAnyOf(ownerName string) func(...string) error {
	return func(fields ...string) error {
		return uf.checkFields(ownerName, fields, false)
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

func CheckTrustedSource(uf *UserFile) error {
	if len(uf.TrustedSource) == 0 {
		return fmt.Errorf(fconsts.MissingTrustedSourceError)
	}
	return nil
}
