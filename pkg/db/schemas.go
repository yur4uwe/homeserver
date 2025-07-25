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
	Size          = "size"
	ContentType   = "content_type"
	OwnerID       = "owner_id"
	UpdatedAt     = "updated_at"
	Metadata      = "metadata"
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

type UserFiles struct {
	ID          int       `db:"id" json:"id"`
	Filename    string    `db:"filename" json:"filename"`
	Path        string    `db:"path" json:"path"`
	Size        int64     `db:"size" json:"size"`
	ContentType string    `db:"content_type" json:"content_type"`
	OwnerID     int       `db:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Metadata    string    `db:"metadata" json:"metadata"`
}
