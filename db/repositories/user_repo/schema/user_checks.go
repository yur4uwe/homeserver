package uschema

import (
	"fmt"
	"homeserver/constants"
	"homeserver/constants/uconsts"
	"slices"
	"strings"
)

func (u *User) ShouldHaveFields(fields ...string) error {
	if len(fields) == 0 {
		return fmt.Errorf(uconsts.NoFieldsSpecifiedError)
	}

	var errorMessages = make(map[string]string)

	for _, field := range fields {
		if !slices.Contains(constants.DBUsersTable.Fields, field) {
			errorMessages[field] = fmt.Sprintf(uconsts.InvalidUserFieldError, field)
			continue
		}

		switch field {
		case uconsts.ID:
			if err := CheckUserID(u); err != nil {
				errorMessages[uconsts.ID] = err.Error()
			}
		case uconsts.Username:
			if err := CheckUsername(u); err != nil {
				errorMessages[uconsts.Username] = err.Error()
			}
		case uconsts.Email:
			if err := CheckEmail(u); err != nil {
				errorMessages[uconsts.Email] = err.Error()
			}
		case uconsts.EmailVerified:
			if err := CheckEmailVerified(u); err != nil {
				errorMessages[uconsts.EmailVerified] = err.Error()
			}
		case uconsts.Password:
			if err := CheckPassword(u); err != nil {
				errorMessages[uconsts.Password] = err.Error()
			}
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("User validation failed: %v", errorMessages)
	}

	return nil
}

func CheckUserID(u *User) error {
	if u.ID <= 0 {
		return fmt.Errorf(uconsts.InvalidIDError)
	}
	return nil
}

func CheckUsername(u *User) error {
	if u.Username == "" {
		return fmt.Errorf(uconsts.MissingUsernameError)
	}
	return nil
}

func CheckEmail(u *User) error {
	if u.Email == "" {
		return fmt.Errorf(uconsts.InvalidEmailError)
	}
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf(uconsts.InvalidEmailFormatError)
	}
	return nil
}

func CheckEmailVerified(u *User) error {
	if !u.EmailVerified {
		return fmt.Errorf(uconsts.InvalidEmailVerifiedError)
	}
	return nil
}

func CheckPassword(u *User) error {
	if u.Password == "" {
		return fmt.Errorf(uconsts.InvalidPasswordError)
	}
	return nil
}
