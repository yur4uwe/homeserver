package folschema

import (
	"fmt"
	"homeserver/constants"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/pkg/utils/generics"
)

func (uf *UserFolder) GetValidFieldsList() []string {
	var validFields []string
	if uf.ID != 0 {
		validFields = append(validFields, folconsts.ID)
	}
	if uf.Name != "" {
		validFields = append(validFields, folconsts.FolderName)
	}
	if uf.Path != "" {
		validFields = append(validFields, folconsts.Path)
	}
	if uf.Parent != nil {
		validFields = append(validFields, folconsts.FolderParent)
	}
	if uf.OwnerID != 0 {
		validFields = append(validFields, folconsts.OwnerID)
	}
	if len(uf.TrustedSource) > 0 {
		validFields = append(validFields, folconsts.TrustedSource)
	}
	if !uf.CreatedAt.IsZero() {
		validFields = append(validFields, folconsts.CreatedAt)
	}

	return validFields
}

func (uf *UserFolder) GetValidFieldsData() []any {
	var validFields []any

	if uf.ID != 0 {
		validFields = append(validFields, uf.ID)
	}
	if uf.Name != "" {
		validFields = append(validFields, uf.Name)
	}
	if uf.Path != "" {
		validFields = append(validFields, uf.Path)
	}
	if uf.Parent != nil {
		validFields = append(validFields, uf.Parent)
	}
	if uf.OwnerID != 0 {
		validFields = append(validFields, uf.OwnerID)
	}
	if len(uf.TrustedSource) > 0 {
		validFields = append(validFields, uf.TrustedSource)
	}
	if !uf.CreatedAt.IsZero() {
		validFields = append(validFields, uf.CreatedAt)
	}

	return validFields
}

func (uf *UserFolder) AbsolutePath() string {
	storagePath := GetStoragePath(constants.StoragePath)

	return fmt.Sprintf("%s%s/%s", storagePath, uf.Path, uf.Name)
}

func FieldsList(without ...string) []string {
	return generics.Without(dbconsts.DBFoldersTable.Fields, without...)
}

func (uf *UserFolder) MergeFolders(other *UserFolder) {
	if other == nil {
		return
	}

	if other.Name != "" {
		uf.Name = other.Name
	}
	if other.Path != "" {
		uf.Path = other.Path
	}
	if other.Parent != nil {
		uf.Parent = other.Parent
	}
	if other.OwnerID != 0 {
		uf.OwnerID = other.OwnerID
	}
	if len(other.TrustedSource) > 0 {
		uf.TrustedSource = other.TrustedSource
	}
}
