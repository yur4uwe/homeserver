package folschema

import (
	"fmt"
	"homeserver/constants"
	"homeserver/web-server/pkg/utils"
	"os"
	"strings"
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
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		panic("Folder name cannot contain slashes")
	}
	if strings.Contains(name, ".") {
		panic("Folder name cannot contain periods")
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
	storagePath := getStoragePath(constants.StoragePath)

	return fmt.Sprintf("%s%s/%s", storagePath, uf.Path, uf.Name)
}

func FieldsList(without ...string) []string {
	return utils.Without(constants.DBFoldersTable.Fields, without...)
}
