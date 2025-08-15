package folder_schema

import "time"

type UserFolder struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	// Path start from user root folder, e.g. "<username>/documents/reports"
	// With no leading or trailing slashes
	Path string `db:"path" json:"path"`
	// If Parent is nil, this is a root user folder
	Parent    *int      `db:"parent" json:"parent,omitempty"`
	OwnerID   int       `db:"owner_id" json:"owner_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
