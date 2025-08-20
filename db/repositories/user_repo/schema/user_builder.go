package uschema

import (
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"homeserver/web-server/pkg/utils"
	"slices"
	"strings"
	"time"
)

var (
	defaultUser = User{
		ID:            0,
		Username:      "",
		Password:      "",
		Email:         "",
		AuthToken:     nil,
		CreatedAt:     time.Time{},
		EmailVerified: false,
	}
)

func NewUserBuilder() *User {
	clone := defaultUser
	return &clone
}

func (u *User) WithID(id int) *User {
	if id <= 0 {
		panic("User ID must be greater than 0")
	}
	u.ID = id
	return u
}

func (u *User) WithUsername(username string) *User {
	if username == "" {
		panic("Username cannot be empty")
	}
	u.Username = username
	return u
}

func (u *User) WithPassword(password string) *User {
	if password == "" {
		panic("Password cannot be empty")
	}
	u.Password = password
	return u
}

func (u *User) WithEmail(email string) *User {
	if email == "" || !strings.Contains(email, "@") {
		panic("Invalid email address")
	}
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

func (u *User) GetValidFields() []string {
	var validFields []string
	for _, field := range FieldsList() {
		switch field {
		case uconsts.ID:
			if u.ID != 0 {
				validFields = append(validFields, field)
			}
		case uconsts.Username:
			if u.Username != "" {
				validFields = append(validFields, field)
			}
		case uconsts.Password:
			if u.Password != "" {
				validFields = append(validFields, field)
			}
		case uconsts.Email:
			if u.Email != "" {
				validFields = append(validFields, field)
			}
		case uconsts.AuthToken:
			if u.AuthToken != nil {
				validFields = append(validFields, field)
			}
		case uconsts.CreatedAt:
			if !u.CreatedAt.IsZero() {
				validFields = append(validFields, field)
			}
		}
	}
	return validFields
}

func (u *User) Get(field string) (any, bool) {
	switch field {
	case uconsts.ID:
		return u.ID, true
	case uconsts.Username:
		return u.Username, true
	case uconsts.Password:
		return u.Password, true
	case uconsts.Email:
		return u.Email, true
	case uconsts.AuthToken:
		return u.AuthToken, true
	case uconsts.CreatedAt:
		return u.CreatedAt, true
	case uconsts.EmailVerified:
		return u.EmailVerified, true
	default:
		return nil, false
	}
}

func IsUserField(field string) bool {
	return slices.Contains(FieldsList(), field)
}

func FieldsList(without ...string) []string {
	return utils.Without(dbconsts.DBUsersTable.Fields, without...)
}
