package schemas

import (
	"fmt"
	"homeserver/pkg/globals"
	"os"
	"time"
)

func NewUserBuilder() *User {
	return &User{
		ID:            0,
		Username:      "",
		Password:      "",
		Email:         "",
		AuthToken:     nil,
		CreatedAt:     time.Now(),
		EmailVerified: false,
	}
}

func (u *User) WithID(id int) *User {
	u.ID = id
	return u
}

func (u *User) WithUsername(username string) *User {
	u.Username = username
	return u
}

func (u *User) WithPassword(password string) *User {
	u.Password = password
	return u
}

func (u *User) WithEmail(email string) *User {
	u.Email = email
	return u
}

func (u *User) WithAuthToken(authToken *string) *User {
	u.AuthToken = authToken
	return u
}

func (u *User) WithEmailVerified(emailVerified bool) *User {
	u.EmailVerified = emailVerified
	return u
}

func NewUserFileBuilder() *UserFile {
	return &UserFile{
		ID:          0,
		Filename:    "",
		Extension:   "",
		Path:        "",
		Directory:   0,
		Size:        0,
		ContentType: "",
		OwnerID:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    nil,
	}
}

func (uf *UserFile) WithID(id int) *UserFile {
	uf.ID = id
	return uf
}

func (uf *UserFile) WithFilename(filename string) *UserFile {
	uf.Filename = filename
	return uf
}

func (uf *UserFile) WithExtension(extension string) *UserFile {
	uf.Extension = extension
	return uf
}

func (uf *UserFile) WithPath(path string) *UserFile {
	uf.Path = path
	return uf
}

func (uf *UserFile) WithDirectory(directory int) *UserFile {
	uf.Directory = directory
	return uf
}

func (uf *UserFile) WithSize(size int64) *UserFile {
	uf.Size = size
	return uf
}

func (uf *UserFile) WithContentType(contentType string) *UserFile {
	uf.ContentType = contentType
	return uf
}

func (uf *UserFile) WithOwnerID(ownerID int) *UserFile {
	uf.OwnerID = ownerID
	return uf
}

func (uf *UserFile) WithMetadata(metadata *string) *UserFile {
	uf.Metadata = metadata
	return uf
}

func (uf *UserFile) AbsolutePath() string {
	return globals.BuildFilePath(uf.Path, uf.Filename, uf.Extension)
}

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
