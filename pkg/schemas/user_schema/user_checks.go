package user_schema

import (
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas"
	"slices"
	"strings"
)

func (u *User) ShouldHaveFields(fields ...string) error {
	if len(fields) == 0 {
		return fmt.Errorf(schemas.NoFieldsSpecifiedError)
	}

	var errorMessages []string

	for _, field := range fields {
		if !slices.Contains(globals.DBUsersTable.Fields, field) {
			errorMessages = append(errorMessages, fmt.Sprintf(schemas.InvalidUserFieldError, field))
			continue
		}

		switch field {
		case schemas.UserID:
			if err := CheckUserID(u); err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("UserID: %v", err))
			}
		case schemas.Username:
			if err := CheckUsername(u); err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("Username: %v", err))
			}
		case schemas.Email:
			if err := CheckEmail(u); err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("Email: %v", err))
			}
		case schemas.EmailVerified:
			if err := CheckEmailVerified(u); err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("EmailVerified: %v", err))
			}
		case schemas.Password:
			if err := CheckPassword(u); err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("Password: %v", err))
			}
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("%s%s", schemas.ValidationErrorPrefix, strings.Join(errorMessages, "\n"))
	}

	return nil
}

func CheckUserID(u *User) error {
	if u.ID <= 0 {
		return fmt.Errorf(schemas.InvalidUserIDError)
	}
	return nil
}

func CheckUsername(u *User) error {
	if u.Username == "" {
		return fmt.Errorf(schemas.InvalidUsernameError)
	}
	return nil
}

func CheckEmail(u *User) error {
	if u.Email == "" {
		return fmt.Errorf(schemas.InvalidEmailError)
	}
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf(schemas.InvalidEmailFormatError)
	}
	return nil
}

func CheckEmailVerified(u *User) error {
	if !u.EmailVerified {
		return fmt.Errorf(schemas.InvalidEmailVerifiedError)
	}
	return nil
}

func CheckPassword(u *User) error {
	if u.Password == "" {
		return fmt.Errorf(schemas.InvalidPasswordError)
	}
	return nil
}
