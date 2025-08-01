package user_schema

import (
	"fmt"
	"strings"
)

func HasID(u *User, err error) (*User, error) {
	if err != nil {
		return u, err
	}
	if u.ID <= 0 {
		return u, fmt.Errorf("user ID is invalid")
	}
	return u, nil
}

func HasValidUsername(u *User, err error) (*User, error) {
	if err != nil {
		return u, err
	}
	if u.Username == "" {
		return u, fmt.Errorf("username is invalid")
	}
	return u, nil
}

func HasValidEmail(u *User, err error) (*User, error) {
	if err != nil {
		return u, err
	}
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		return u, fmt.Errorf("email is invalid")
	}
	return u, nil
}

func HasEmailVerified(u *User, err error) (*User, error) {
	if err != nil {
		return u, err
	}
	if !u.EmailVerified {
		return u, fmt.Errorf("email is not verified")
	}
	return u, nil
}

func HasValidPassword(u *User, err error) (*User, error) {
	if err != nil {
		return u, err
	}
	return u, nil
}
