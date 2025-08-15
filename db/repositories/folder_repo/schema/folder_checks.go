package folder_schema

import (
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas"
	"slices"
	"strings"
)

func (uf *UserFolder) ShouldHaveFields(ownerName string) func(...string) error {
	return func(fields ...string) error {
		if len(fields) == 0 {
			return fmt.Errorf(schemas.NoFieldsSpecifiedError)
		}

		var errorMessages []string

		for _, field := range fields {
			if !slices.Contains(globals.DBFoldersTable.Fields, field) {
				errorMessages = append(errorMessages, fmt.Sprintf(schemas.InvalidFolderFieldError, field))
				continue
			}

			switch field {
			case schemas.FolderID:
				if err := CheckFolderID(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("FolderID: %v", err))
				}
			case schemas.FolderName:
				if err := CheckFolderName(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("FolderName: %v", err))
				}
			case schemas.Path:
				if err := CheckFolderPath(ownerName)(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Path: %v", err))
				}
			case schemas.FolderParent:
				if err := CheckParent(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("FolderParent: %v", err))
				}
			case schemas.OwnerID:
				if err := CheckFolderOwnerID(uf); err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("OwnerID: %v", err))
				}
			}
		}

		if len(errorMessages) > 0 {
			return fmt.Errorf("%s%s", schemas.ValidationErrorPrefix, strings.Join(errorMessages, "\n"))
		}

		return nil
	}
}

func CheckFolderID(uf *UserFolder) error {
	if uf.ID == 0 {
		return fmt.Errorf(schemas.InvalidFolderIDError)
	}
	return nil
}

func CheckFolderName(uf *UserFolder) error {
	if uf.Name == "" {
		return fmt.Errorf(schemas.InvalidFolderNameError)
	}
	return nil
}

func CheckFolderPath(ownerName string) func(uf *UserFolder) error {
	return func(uf *UserFolder) error {
		if uf.Path == "" {
			return fmt.Errorf(schemas.InvalidFolderPathError)
		}
		if strings.HasPrefix(uf.Path, "/") {
			return fmt.Errorf(schemas.InvalidFolderPathStartError)
		}
		if strings.HasSuffix(uf.Path, "/") {
			return fmt.Errorf(schemas.InvalidFolderPathEndError)
		}
		if !strings.HasPrefix(uf.Path, ownerName) {
			return fmt.Errorf(schemas.InvalidFolderPathOwnerError, ownerName)
		}
		return nil
	}
}

func CheckParent(uf *UserFolder) error {
	if uf.Parent != nil && *uf.Parent <= 0 {
		return fmt.Errorf(schemas.InvalidFolderParentError)
	}
	return nil
}

func CheckFolderOwnerID(uf *UserFolder) error {
	if uf.OwnerID == 0 {
		return fmt.Errorf(schemas.InvalidFolderOwnerIDError)
	}
	return nil
}
