package folder_schema

import (
	"fmt"
	"strings"
)

func CheckFolderID(uf *UserFolder, err error) (*UserFolder, error) {
	if err != nil {
		return uf, err
	}
	if uf.ID == 0 {
		return uf, fmt.Errorf("folder ID is invalid: ID must be greater than 0")
	}
	return uf, nil
}

func CheckFolderName(uf *UserFolder, err error) (*UserFolder, error) {
	if err != nil {
		return uf, err
	}
	if uf.Name == "" {
		return uf, fmt.Errorf("folder name is invalid: name cannot be empty")
	}
	return uf, nil
}

func CheckFolderPath(ownerName string) func(uf *UserFolder, err error) (*UserFolder, error) {
	return func(uf *UserFolder, err error) (*UserFolder, error) {
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

func CheckParent(uf *UserFolder, err error) (*UserFolder, error) {
	if err != nil {
		return uf, err
	}
	if uf.Parent != nil && *uf.Parent <= 0 {
		return uf, fmt.Errorf("parent is invalid: parent ID must be greater than 0")
	}
	return uf, nil
}

func CheckFolderOwnerID(uf *UserFolder, err error) (*UserFolder, error) {
	if err != nil {
		return uf, err
	}
	if uf.OwnerID == 0 {
		return uf, fmt.Errorf("owner ID is invalid: owner ID must be greater than 0")
	}
	return uf, nil
}
