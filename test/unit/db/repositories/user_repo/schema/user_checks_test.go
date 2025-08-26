package uschema

import (
	"homeserver/constants/uconsts"
	uschema "homeserver/db/repositories/user_repo/schema"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserChecksSuite(t *testing.T) {
	t.Run("TestCheckUserID", func(t *testing.T) {
		user := &uschema.User{ID: 1}
		err := uschema.CheckUserID(user)
		assert.NoError(t, err, "CheckUserID should not return an error for valid ID")

		user = &uschema.User{ID: 0}
		err = uschema.CheckUserID(user)
		assert.Error(t, err, "CheckUserID should return an error for invalid ID")
		assert.EqualError(t, err, "ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckUsername", func(t *testing.T) {
		user := &uschema.User{Username: "john_doe"}
		err := uschema.CheckUsername(user)
		assert.NoError(t, err, "CheckUsername should not return an error for valid username")

		user = &uschema.User{Username: ""}
		err = uschema.CheckUsername(user)
		assert.Error(t, err, "CheckUsername should return an error for empty username")
		assert.EqualError(t, err, "username cannot be empty", "Error message should match")
	})

	t.Run("TestCheckEmail", func(t *testing.T) {
		user := &uschema.User{Email: "john@example.com"}
		err := uschema.CheckEmail(user)
		assert.NoError(t, err, "CheckEmail should not return an error for valid email")

		user = &uschema.User{Email: ""}
		err = uschema.CheckEmail(user)
		assert.Error(t, err, "CheckEmail should return an error for empty email")
		assert.EqualError(t, err, "email cannot be empty", "Error message should match")

		user = &uschema.User{Email: "invalid_email"}
		err = uschema.CheckEmail(user)
		assert.Error(t, err, "CheckEmail should return an error for invalid email")
		assert.EqualError(t, err, "email must contain '@'", "Error message should match")
	})

	t.Run("TestCheckEmailVerified", func(t *testing.T) {
		user := &uschema.User{EmailVerified: true}
		err := uschema.CheckEmailVerified(user)
		assert.NoError(t, err, "CheckEmailVerified should not return an error for verified email")

		user = &uschema.User{EmailVerified: false}
		err = uschema.CheckEmailVerified(user)
		assert.Error(t, err, "CheckEmailVerified should return an error for unverified email")
		assert.EqualError(t, err, "email is not verified", "Error message should match")
	})

	t.Run("TestCheckPassword", func(t *testing.T) {
		user := &uschema.User{Password: "secure_password"}
		err := uschema.CheckPassword(user)
		assert.NoError(t, err, "CheckPassword should not return an error for valid password")

		user = &uschema.User{Password: ""}
		err = uschema.CheckPassword(user)
		assert.Error(t, err, "CheckPassword should return an error for empty password")
		assert.EqualError(t, err, "password cannot be empty", "Error message should match")
	})

	t.Run("TestShouldHaveFields", func(t *testing.T) {
		user := &uschema.User{
			ID:            0,
			Username:      "",
			Email:         "invalid_email",
			EmailVerified: false,
			Password:      "",
		}

		err := user.ShouldHaveFields(uconsts.ID, uconsts.Username, uconsts.Email, uconsts.EmailVerified, uconsts.Password)
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid fields")
		assert.Contains(t, err.Error(), uconsts.InvalidIDError)
		assert.Contains(t, err.Error(), uconsts.MissingUsernameError)
		assert.Contains(t, err.Error(), uconsts.InvalidEmailFormatError)
		assert.Contains(t, err.Error(), uconsts.InvalidEmailVerifiedError)
		assert.Contains(t, err.Error(), uconsts.InvalidPasswordError)

		user = &uschema.User{
			ID:            1,
			Username:      "john_doe",
			Email:         "john@example.com",
			EmailVerified: true,
			Password:      "secure_password",
		}

		err = user.ShouldHaveFields(uconsts.ID, uconsts.Username, uconsts.Email, uconsts.EmailVerified, uconsts.Password)
		assert.NoError(t, err, "ShouldHaveFields should not return an error for valid fields")

		err = user.ShouldHaveFields()
		assert.Error(t, err, "ShouldHaveFields should return an error when no fields are specified")

		err = user.ShouldHaveFields("invalid_field")
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid field")
	})
}
