package folder_schema

import (
	"fmt"
	"homeserver/pkg/globals"
	"os"
	"time"
)

var (
	getStoragePath = os.Getenv
)

var (
	defaultFolder = UserFolder{
		ID:        0,
		Name:      "",
		Path:      "",
		Parent:    nil,
		OwnerID:   0,
		CreatedAt: time.Now(),
	}
)

func NewUserFolderBuilder() *UserFolder {
	clone := defaultFolder
	return &clone
}

func (uf *UserFolder) WithID(id int) *UserFolder {
	if id <= 0 {
		panic("Folder ID must be greater than 0")
	}
	uf.ID = id
	return uf
}

func (uf *UserFolder) WithName(name string) *UserFolder {
	if name == "" {
		panic("Folder name cannot be empty")
	}
	uf.Name = name
	return uf
}

func (uf *UserFolder) WithPath(path string) *UserFolder {
	if path == "" {
		panic("Folder path cannot be empty")
	}
	if path[0] == '/' || path[len(path)-1] == '/' {
		panic("Folder path cannot start or end with a slash")
	}
	uf.Path = path
	return uf
}

func (uf *UserFolder) WithParent(parent *int) *UserFolder {
	if parent != nil && *parent <= 0 {
		panic("Parent folder ID must be greater than 0")
	}
	uf.Parent = parent
	return uf
}

func (uf *UserFolder) WithOwnerID(ownerID int) *UserFolder {
	if ownerID <= 0 {
		panic("Owner ID must be greater than 0")
	}
	uf.OwnerID = ownerID
	return uf
}

func (uf *UserFolder) AbsolutePath() string {
	storagePath := getStoragePath(globals.EnvStoragePath)

	return fmt.Sprintf("%s%s/%s", storagePath, uf.Path, uf.Name)
}
