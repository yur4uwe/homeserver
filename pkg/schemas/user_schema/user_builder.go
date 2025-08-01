package user_schema

import (
	"strings"
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
