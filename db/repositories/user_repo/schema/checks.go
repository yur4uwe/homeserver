package uschema

import (
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"slices"
	"strings"
)

func (u *User) checkFields(fields []string, should_have bool) error {
	if len(fields) == 0 {
		return fmt.Errorf(uconsts.NoFieldsSpecifiedError)
	}

	var errorMessages []string

	for _, field := range fields {
		if !slices.Contains(dbconsts.DBUsersTable.Fields, field) {
			errorMessages = append(errorMessages, fmt.Sprintf(uconsts.InvalidUserFieldError, field))
			continue
		}

		var err error
		switch field {
		case uconsts.ID:
			err = CheckUserID(u)
		case uconsts.Username:
			err = CheckUsername(u)
		case uconsts.Password:
			err = CheckPassword(u)
		case uconsts.AuthToken:
			err = CheckAuthToken(u)
		case uconsts.CreatedAt:
			err = CheckCreatedAt(u)
		case uconsts.Email:
			err = CheckEmail(u)

		}

		if should_have {
			if err != nil {
				errorMessages = append(errorMessages, err.Error())
			}
		} else {
			if err == nil {
				errorMessages = append(errorMessages, fmt.Sprintf("%s should not be set", field))
			}
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("User validation failed: %v", errorMessages)
	}

	return nil
}

func (u *User) ShouldNotHaveFields(fields ...string) error {
	return u.checkFields(fields, false)
}

func (u *User) ShouldHaveFields(fields ...string) error {
	return u.checkFields(fields, true)
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

func CheckCreatedAt(u *User) error {
	if u.CreatedAt.IsZero() {
		return fmt.Errorf(uconsts.InvalidCreatedAtError)
	}
	return nil
}

func CheckAuthToken(u *User) error {
	if u.AuthToken != nil && *u.AuthToken == "" {
		return fmt.Errorf(uconsts.InvalidAuthTokenError)
	}
	return nil
}
