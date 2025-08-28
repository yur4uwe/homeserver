package fschema

import (
	"fmt"
	"homeserver/constants"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"homeserver/pkg/utils/generics"
	"os"
	"time"
)

var (
	Stat = os.Stat
)

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
	return generics.Without(dbconsts.DBFilesTable.Fields, without...)
}

func (uf *UserFile) GetValidFieldsData() []any {
	var validFields []any
	if uf.ID != 0 {
		validFields = append(validFields, uf.ID)
	}
	if uf.Filename != "" {
		validFields = append(validFields, uf.Filename)
	}
	if uf.Extension != "" {
		validFields = append(validFields, uf.Extension)
	}
	if uf.Path != "" {
		validFields = append(validFields, uf.Path)
	}
	if uf.Directory != 0 {
		validFields = append(validFields, uf.Directory)
	}
	if uf.Size != 0 {
		validFields = append(validFields, uf.Size)
	}
	if uf.ContentType != "" {
		validFields = append(validFields, uf.ContentType)
	}
	if uf.OwnerID != 0 {
		validFields = append(validFields, uf.OwnerID)
	}
	if !uf.CreatedAt.IsZero() {
		validFields = append(validFields, uf.CreatedAt)
	}
	if !uf.UpdatedAt.IsZero() {
		validFields = append(validFields, uf.UpdatedAt)
	}
	return validFields
}

func (uf *UserFile) GetValidFieldsList() []string {
	var validFields []string
	if uf.ID > 0 {
		validFields = append(validFields, fconsts.ID)
	}
	if uf.Filename != "" {
		validFields = append(validFields, fconsts.Filename)
	}
	if uf.Extension != "" {
		validFields = append(validFields, fconsts.Extension)
	}
	if uf.Path != "" {
		validFields = append(validFields, fconsts.Path)
	}
	if uf.Directory > 0 {
		validFields = append(validFields, fconsts.Directory)
	}
	if uf.Size > 0 {
		validFields = append(validFields, fconsts.Size)
	}
	if uf.ContentType != "" {
		validFields = append(validFields, fconsts.ContentType)
	}
	if uf.OwnerID > 0 {
		validFields = append(validFields, fconsts.OwnerID)
	}
	return validFields
}

func (uf *UserFile) MergeFiles(new_file *UserFile) {
	if new_file == nil {
		return
	}

	if new_file.Filename != "" {
		uf.Filename = new_file.Filename
	}

	if new_file.Extension != "" {
		uf.Extension = new_file.Extension
	}

	if new_file.Path != "" {
		uf.Path = new_file.Path
	}

	if new_file.Directory != 0 {
		uf.Directory = new_file.Directory
	}

	if new_file.Size != 0 {
		uf.Size = new_file.Size
	}

	if new_file.ContentType != "" {
		uf.ContentType = new_file.ContentType
	}

	if new_file.Metadata != nil {
		uf.Metadata = new_file.Metadata
	}

	uf.UpdatedAt = time.Now()

	new_file, err := PreventDuplicateFileName(uf)
	if err != nil {
		panic(err)
	}

	*uf = *new_file
}

func PreventDuplicateFileName(file *UserFile) (*UserFile, error) {
	new_file := *file
	counter := 1
	for {
		path := new_file.AbsolutePath()

		if _, err := Stat(path); err != nil && os.IsNotExist(err) {
			return &new_file, nil
		}
		fmt.Printf("File %s already exists, trying to prevent duplicate...\n", new_file.AbsolutePath())

		new_file.Filename = fmt.Sprintf("%s_%d", new_file.Filename, counter)

		counter++
	}
}
