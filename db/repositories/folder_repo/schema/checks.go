package folschema

import (
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"slices"
	"strings"
)

func (uf *UserFolder) checkFields(ownerName string, fields []string, shouldHave bool) error {
	if len(fields) == 0 {
		return fmt.Errorf(folconsts.NoFieldsSpecifiedError)
	}

	var errorMessages []string

	for _, field := range fields {
		if !slices.Contains(dbconsts.DBFoldersTable.Fields, field) {
			errorMessages = append(errorMessages, fmt.Sprintf(folconsts.InvalidFolderFieldError, field))
			continue
		}

		var err error
		switch field {
		case folconsts.ID:
			err = CheckFolderID(uf)
		case folconsts.FolderName:
			err = CheckFolderName(uf)
		case folconsts.Path:
			err = CheckFolderPath(ownerName)(uf)
		case folconsts.FolderParent:
			err = CheckParent(uf)
		case folconsts.OwnerID:
			err = CheckFolderOwnerID(uf)
		case folconsts.TrustedSource:
			err = CheckFolderTrustedSource(uf)
		case folconsts.CreatedAt:
			if uf.CreatedAt.IsZero() {
				err = fmt.Errorf(folconsts.InvalidCreatedAtError)
			}
		default:
			err = fmt.Errorf(folconsts.InvalidFolderFieldError, field)
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

func (uf *UserFolder) ShouldHaveFields(ownerName string) func(...string) error {
	return func(fields ...string) error {
		return uf.checkFields(ownerName, fields, true)
	}
}

func (uf *UserFolder) ShouldNotHaveFields(ownerName string) func(...string) error {
	return func(fields ...string) error {
		return uf.checkFields(ownerName, fields, false)
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

func CheckFolderTrustedSource(uf *UserFolder) error {
	if len(uf.TrustedSource) == 0 {
		return fmt.Errorf(folconsts.InvalidFolderTrustedSourceError)
	}
	return nil
}
