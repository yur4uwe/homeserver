package user_schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserChecksSuite(t *testing.T) {
	t.Run("TestHasID", func(t *testing.T) {
		user := &User{ID: 1}
		result, err := HasID(user, nil)
		assert.NoError(t, err, "HasID should not return an error for valid ID")
		assert.Equal(t, user, result, "HasID should return the same user object")

		user = &User{ID: 0}
		_, err = HasID(user, nil)
		assert.Error(t, err, "HasID should return an error for invalid ID")
		assert.EqualError(t, err, "user ID is invalid", "Error message should match")
	})

	t.Run("TestHasValidUsername", func(t *testing.T) {
		user := &User{Username: "john_doe"}
		result, err := HasValidUsername(user, nil)
		assert.NoError(t, err, "HasValidUsername should not return an error for valid username")
		assert.Equal(t, user, result, "HasValidUsername should return the same user object")

		user = &User{Username: ""}
		_, err = HasValidUsername(user, nil)
		assert.Error(t, err, "HasValidUsername should return an error for empty username")
		assert.EqualError(t, err, "username is invalid", "Error message should match")
	})

	t.Run("TestHasValidEmail", func(t *testing.T) {
		user := &User{Email: "john@example.com"}
		result, err := HasValidEmail(user, nil)
		assert.NoError(t, err, "HasValidEmail should not return an error for valid email")
		assert.Equal(t, user, result, "HasValidEmail should return the same user object")

		user = &User{Email: ""}
		_, err = HasValidEmail(user, nil)
		assert.Error(t, err, "HasValidEmail should return an error for empty email")
		assert.EqualError(t, err, "email is invalid", "Error message should match")

		user = &User{Email: "invalid_email"}
		_, err = HasValidEmail(user, nil)
		assert.Error(t, err, "HasValidEmail should return an error for invalid email")
		assert.EqualError(t, err, "email is invalid", "Error message should match")
	})

	t.Run("TestHasEmailVerified", func(t *testing.T) {
		user := &User{EmailVerified: true}
		result, err := HasEmailVerified(user, nil)
		assert.NoError(t, err, "HasEmailVerified should not return an error for verified email")
		assert.Equal(t, user, result, "HasEmailVerified should return the same user object")

		user = &User{EmailVerified: false}
		_, err = HasEmailVerified(user, nil)
		assert.Error(t, err, "HasEmailVerified should return an error for unverified email")
		assert.EqualError(t, err, "email is not verified", "Error message should match")
	})

	t.Run("TestHasValidPassword", func(t *testing.T) {
		user := &User{Password: "secure_password"}
		result, err := HasValidPassword(user, nil)
		assert.NoError(t, err, "HasValidPassword should not return an error for valid password")
		assert.Equal(t, user, result, "HasValidPassword should return the same user object")

		user = &User{Password: ""}
		result, err = HasValidPassword(user, nil)
		assert.NoError(t, err, "HasValidPassword does not currently validate password presence")
		assert.Equal(t, user, result, "HasValidPassword should return the same user object")
	})
}
