package user_schema

import (
	"fmt"
	"homeserver/pkg/schemas"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserChecksSuite(t *testing.T) {
	t.Run("TestCheckUserID", func(t *testing.T) {
		user := &User{ID: 1}
		err := CheckUserID(user)
		assert.NoError(t, err, "CheckUserID should not return an error for valid ID")

		user = &User{ID: 0}
		err = CheckUserID(user)
		assert.Error(t, err, "CheckUserID should return an error for invalid ID")
		assert.EqualError(t, err, "user ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckUsername", func(t *testing.T) {
		user := &User{Username: "john_doe"}
		err := CheckUsername(user)
		assert.NoError(t, err, "CheckUsername should not return an error for valid username")

		user = &User{Username: ""}
		err = CheckUsername(user)
		assert.Error(t, err, "CheckUsername should return an error for empty username")
		assert.EqualError(t, err, "username cannot be empty", "Error message should match")
	})

	t.Run("TestCheckEmail", func(t *testing.T) {
		user := &User{Email: "john@example.com"}
		err := CheckEmail(user)
		assert.NoError(t, err, "CheckEmail should not return an error for valid email")

		user = &User{Email: ""}
		err = CheckEmail(user)
		assert.Error(t, err, "CheckEmail should return an error for empty email")
		assert.EqualError(t, err, "email cannot be empty", "Error message should match")

		user = &User{Email: "invalid_email"}
		err = CheckEmail(user)
		assert.Error(t, err, "CheckEmail should return an error for invalid email")
		assert.EqualError(t, err, "email must contain '@'", "Error message should match")
	})

	t.Run("TestCheckEmailVerified", func(t *testing.T) {
		user := &User{EmailVerified: true}
		err := CheckEmailVerified(user)
		assert.NoError(t, err, "CheckEmailVerified should not return an error for verified email")

		user = &User{EmailVerified: false}
		err = CheckEmailVerified(user)
		assert.Error(t, err, "CheckEmailVerified should return an error for unverified email")
		assert.EqualError(t, err, "email is not verified", "Error message should match")
	})

	t.Run("TestCheckPassword", func(t *testing.T) {
		user := &User{Password: "secure_password"}
		err := CheckPassword(user)
		assert.NoError(t, err, "CheckPassword should not return an error for valid password")

		user = &User{Password: ""}
		err = CheckPassword(user)
		assert.Error(t, err, "CheckPassword should return an error for empty password")
		assert.EqualError(t, err, "password cannot be empty", "Error message should match")
	})

	t.Run("TestShouldHaveFields", func(t *testing.T) {
		user := &User{
			ID:            0,
			Username:      "",
			Email:         "invalid_email",
			EmailVerified: false,
			Password:      "",
		}

		err := user.ShouldHaveFields(schemas.UserID, schemas.Username, schemas.Email, schemas.EmailVerified, schemas.Password)
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid fields")
		assert.EqualError(t, err,
			fmt.Sprintf("%sUserID: %s\nUsername: %s\nEmail: %s\nEmailVerified: %s\nPassword: %s",
				schemas.ValidationErrorPrefix,
				schemas.InvalidUserIDError,
				schemas.InvalidUsernameError,
				schemas.InvalidEmailFormatError,
				schemas.InvalidEmailVerifiedError,
				schemas.InvalidPasswordError,
			),
			"Error message should match",
		)

		user = &User{
			ID:            1,
			Username:      "john_doe",
			Email:         "john@example.com",
			EmailVerified: true,
			Password:      "secure_password",
		}

		err = user.ShouldHaveFields(schemas.UserID, schemas.Username, schemas.Email, schemas.EmailVerified, schemas.Password)
		assert.NoError(t, err, "ShouldHaveFields should not return an error for valid fields")

		err = user.ShouldHaveFields()
		assert.Error(t, err, "ShouldHaveFields should return an error when no fields are specified")

		err = user.ShouldHaveFields("invalid_field")
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid field")
	})
}
