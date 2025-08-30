package uschema

import "time"

type User struct {
	ID            int       `db:"id" json:"id"`
	Username      string    `db:"username" json:"username"`
	Password      string    `db:"password" json:"password"`
	Email         string    `db:"email" json:"email"`
	AuthToken     *string   `db:"auth_token" json:"auth_token,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	EmailVerified bool      `db:"email_verified" json:"email_verified"`
}
