package schemas

import (
	"strings"
)

func (u *User) HasID() *User {
	if u.ID == 0 {
		panic("User ID is invalid")
	}
	return u
}

func (u *User) HasValidUsername() *User {
	if u.Username == "" {
		panic("Username is invalid")
	}
	return u
}

func (u *User) HasValidPassword() *User {
	if u.Password == "" {
		panic("Password is invalid")
	}
	return u
}

func (u *User) HasValidEmail() *User {
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		panic("Email is invalid")
	}
	return u
}

func (u *User) IsEmailVerified() *User {
	if !u.EmailVerified {
		panic("Email is not verified")
	}
	return u
}

func (uf *UserFile) HasID() *UserFile {
	if uf.ID == 0 {
		panic("File ID is invalid")
	}
	return uf
}

func (uf *UserFile) HasValidFilename() *UserFile {
	if uf.Filename == "" {
		panic("Filename is invalid")
	}
	return uf
}

func (uf *UserFile) HasValidExtension() *UserFile {
	if uf.Extension == "" || strings.HasPrefix(uf.Extension, ".") {
		panic("Extension is invalid")
	}
	return uf
}

func (uf *UserFile) HasValidPath(ownerName string) *UserFile {
	if uf.Path == "" || strings.HasPrefix(uf.Path, "/") || strings.HasSuffix(uf.Path, "/") || !strings.HasPrefix(uf.Path, ownerName) {
		panic("Path is invalid")
	}
	return uf
}

func (uf *UserFile) HasValidDirectory() *UserFile {
	if uf.Directory == 0 {
		panic("Directory is invalid")
	}
	return uf
}

func (uf *UserFile) HasValidSize() *UserFile {
	if uf.Size <= 0 {
		panic("Size is invalid")
	}
	return uf
}

func (uf *UserFile) HasValidContentType() *UserFile {
	if uf.ContentType == "" {
		panic("Content type is invalid")
	}
	return uf
}

func (uf *UserFile) HasValidOwnerID() *UserFile {
	if uf.OwnerID == 0 {
		panic("Owner ID is invalid")
	}
	return uf
}

func (uf *UserFile) HasMetadata() *UserFile {
	if uf.Metadata == nil || *uf.Metadata == "" {
		panic("Metadata is invalid")
	}
	return uf
}

func (uf *UserFolder) HasID() *UserFolder {
	if uf.ID == 0 {
		panic("Folder ID is invalid")
	}
	return uf
}

func (uf *UserFolder) HasValidName() *UserFolder {
	if uf.Name == "" {
		panic("Folder name is invalid")
	}
	return uf
}

func (uf *UserFolder) HasValidPath(ownerName string) *UserFolder {
	if uf.Path == "" || strings.HasPrefix(uf.Path, "/") || strings.HasSuffix(uf.Path, "/") || !strings.HasPrefix(uf.Path, ownerName) {
		panic("Path is invalid")
	}
	return uf
}

func (uf *UserFolder) HasValidParent() *UserFolder {
	if uf.Parent != nil && *uf.Parent <= 0 {
		panic("Parent is invalid")
	}
	return uf
}

func (uf *UserFolder) HasValidOwnerID() *UserFolder {
	if uf.OwnerID == 0 {
		panic("Owner ID is invalid")
	}
	return uf
}
