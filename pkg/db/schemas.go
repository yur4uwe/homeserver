package db

import "time"

const (
	UserID        = "id"
	Username      = "username"
	Password      = "password"
	Email         = "email"
	AuthToken     = "auth_token"
	EmailVerified = "email_verified"
	CreatedAt     = "created_at"
	FileID        = "id"
	Filename      = "filename"
	Path          = "path"
	Directory     = "directory_id"
	Size          = "size"
	ContentType   = "content_type"
	OwnerID       = "owner_id"
	UpdatedAt     = "updated_at"
	Metadata      = "metadata"
	FolderID      = "id"
	FolderName    = "name"
	FolderParent  = "parent"
)

type User struct {
	ID            int       `db:"id" json:"id"`
	Username      string    `db:"username" json:"username"`
	Password      string    `db:"password" json:"password"`
	Email         string    `db:"email" json:"email"`
	AuthToken     *string   `db:"auth_token" json:"auth_token,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	EmailVerified bool      `db:"email_verified" json:"email_verified"`
}

type UserFile struct {
	ID       int    `db:"id" json:"id"`
	Filename string `db:"filename" json:"filename"`
	// Path starts from user root folder, e.g. "root/documents/reports"
	// With no leading or trailing slashes
	Path string `db:"path" json:"path"`
	// Directory is the ID of the folder this file belongs to
	Directory   int       `db:"directory_id" json:"directory_id"`
	Size        int64     `db:"size" json:"size"`
	ContentType string    `db:"content_type" json:"content_type"`
	OwnerID     int       `db:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Metadata    *string   `db:"metadata" json:"metadata"`
}

type UserFolder struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	// Path start from user root folder, e.g. "root/documents/reports"
	// With no leading or trailing slashes
	Path string `db:"path" json:"path"`
	// If Parent is nil, this is a root folder for the user
	Parent    *int      `db:"parent" json:"parent,omitempty"`
	OwnerID   int       `db:"owner_id" json:"owner_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
