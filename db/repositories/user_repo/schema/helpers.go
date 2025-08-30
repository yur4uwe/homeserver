package uschema

import (
	"fmt"
	"homeserver/constants/uconsts"
)

func (u *User) Identifier(placeholder int) (string, []any) {
	id_err := u.ShouldHaveFields(uconsts.ID)
	uname_err := u.ShouldHaveFields(uconsts.Username)
	email_err := u.ShouldHaveFields(uconsts.Email)

	args := make([]any, 0, 1)
	whereClause := ""
	if id_err == nil {
		whereClause = fmt.Sprintf("id = $%d", placeholder)
		args = append(args, u.ID)
	} else if uname_err == nil {
		whereClause = fmt.Sprintf("username = $%d", placeholder)
		args = append(args, u.Username)
	} else if email_err == nil {
		whereClause = fmt.Sprintf("email = $%d", placeholder)
		args = append(args, u.Email)
	}

	return whereClause, args
}

func (u *User) Merge(another *User) {
	if another == nil {
		return
	}

	if another.ID != 0 {
		u.ID = another.ID
	}
	if another.Username != "" {
		u.Username = another.Username
	}
	if another.Email != "" {
		u.Email = another.Email
	}
	if another.Password != "" {
		u.Password = another.Password
	}
	if !another.CreatedAt.IsZero() {
		u.CreatedAt = another.CreatedAt
	}
	if another.AuthToken != nil {
		u.AuthToken = another.AuthToken
	}
}
