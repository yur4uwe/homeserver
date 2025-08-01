package folder_schema

import (
	"fmt"
	"homeserver/pkg/globals"
	"os"
	"time"
)

func NewUserFolderBuilder() *UserFolder {
	return &UserFolder{
		ID:        0,
		Name:      "",
		Path:      "",
		Parent:    nil,
		OwnerID:   0,
		CreatedAt: time.Now(),
	}
}

func (uf *UserFolder) WithID(id int) *UserFolder {
	uf.ID = id
	return uf
}

func (uf *UserFolder) WithName(name string) *UserFolder {
	uf.Name = name
	return uf
}

func (uf *UserFolder) WithPath(path string) *UserFolder {
	uf.Path = path
	return uf
}

func (uf *UserFolder) WithParent(parent *int) *UserFolder {
	uf.Parent = parent
	return uf
}

func (uf *UserFolder) WithOwnerID(ownerID int) *UserFolder {
	uf.OwnerID = ownerID
	return uf
}

func (uf *UserFolder) AbsolutePath() string {
	storage_path := os.Getenv(globals.EnvStoragePath)

	return fmt.Sprintf("%s%s/%s", storage_path, uf.Path, uf.Name)
}
