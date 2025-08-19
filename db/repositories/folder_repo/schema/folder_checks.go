package folschema

import (
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"slices"
	"strings"
)

func (uf *UserFolder) ShouldHaveFields(ownerName string) func(...string) error {
	return func(fields ...string) error {
		if len(fields) == 0 {
			return fmt.Errorf(folconsts.NoFieldsSpecifiedError)
		}

		var errorMessages = make(map[string]string)

		for _, field := range fields {
			if !slices.Contains(dbconsts.DBFoldersTable.Fields, field) {
				errorMessages[field] = fmt.Sprintf(folconsts.InvalidFolderFieldError, field)
				continue
			}

			switch field {
			case folconsts.ID:
				if err := CheckFolderID(uf); err != nil {
					errorMessages[field] = err.Error()
				}
			case folconsts.FolderName:
				if err := CheckFolderName(uf); err != nil {
					errorMessages[field] = err.Error()
				}
			case folconsts.Path:
				if err := CheckFolderPath(ownerName)(uf); err != nil {
					errorMessages[field] = err.Error()
				}
			case folconsts.FolderParent:
				if err := CheckParent(uf); err != nil {
					errorMessages[field] = err.Error()
				}
			case folconsts.OwnerID:
				if err := CheckFolderOwnerID(uf); err != nil {
					errorMessages[field] = err.Error()
				}
			}
		}

		if len(errorMessages) > 0 {
			return fmt.Errorf("folder validation failed: %v", errorMessages)
		}

		return nil
	}
}

func CheckFolderID(uf *UserFolder) error {
	if uf.ID == 0 {
		return fmt.Errorf(folconsts.InvalidIDError)
	}
	return nil
}

func CheckFolderName(uf *UserFolder) error {
	if uf.Name == "" {
		return fmt.Errorf(folconsts.MissingNameError)
	}
	return nil
}

func CheckFolderPath(ownerName string) func(uf *UserFolder) error {
	return func(uf *UserFolder) error {
		if uf.Path == "" {
			return fmt.Errorf(folconsts.MissingPathError)
		}
		if strings.HasPrefix(uf.Path, "/") {
			return fmt.Errorf(folconsts.InvalidPathStartError)
		}
		if strings.HasSuffix(uf.Path, "/") {
			return fmt.Errorf(folconsts.InvalidPathEndError)
		}
		if !strings.HasPrefix(uf.Path, ownerName) {
			return fmt.Errorf(folconsts.MissingOwnerFolderError)
		}
		return nil
	}
}

func CheckParent(uf *UserFolder) error {
	if uf.Parent != nil && *uf.Parent <= 0 {
		return fmt.Errorf(folconsts.InvalidFolderParentError)
	}
	return nil
}

func CheckFolderOwnerID(uf *UserFolder) error {
	if uf.OwnerID == 0 {
		return fmt.Errorf(folconsts.InvalidOwnerIDError)
	}
	return nil
}
