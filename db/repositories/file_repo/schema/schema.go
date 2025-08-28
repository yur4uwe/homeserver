package fschema

import "time"

type UserFile struct {
	ID        int    `db:"id" json:"id"`
	Filename  string `db:"filename" json:"filename"`
	Extension string `db:"extension" json:"extension"`
	// Path starts from user root folder, e.g. "<username>/documents/reports"
	// With no leading or trailing slashes
	Path string `db:"path" json:"path"`
	// Directory is the ID of the folder this file belongs to
	Directory     int       `db:"directory_id" json:"directory_id"`
	OwnerID       int       `db:"owner_id" json:"owner_id"`
	Size          int64     `db:"size" json:"size"`
	ContentType   string    `db:"content_type" json:"content_type"`
	TrustedSource []byte    `db:"trusted_source" json:"trusted_source"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	Metadata      *string   `db:"metadata" json:"metadata"`
}
