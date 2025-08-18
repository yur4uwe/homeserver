package uschema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserBuilderSuite(t *testing.T) {
	t.Run("TestNewUserBuilder", func(t *testing.T) {
		new_user := NewUserBuilder()

		assert.NotNil(t, new_user, "New user should not be nil")
		assert.Equal(t, defaultUser, *new_user, "New user should match the default user")
	})

	t.Run("TestWithID", func(t *testing.T) {
		new_user := NewUserBuilder().WithID(1)
		assert.Equal(t, 1, new_user.ID, "User ID should be set to 1")
	})

	t.Run("TestWithIDPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserBuilder().WithID(0)
		}, "Setting user ID to 0 should panic")
	})

	t.Run("TestWithUsername", func(t *testing.T) {
		new_user := NewUserBuilder().WithUsername("john_doe")
		assert.Equal(t, "john_doe", new_user.Username, "Username should be set to 'john_doe'")
	})

	t.Run("TestWithUsernamePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserBuilder().WithUsername("")
		}, "Setting an empty username should panic")
	})

	t.Run("TestWithPassword", func(t *testing.T) {
		new_user := NewUserBuilder().WithPassword("secure_password")
		assert.Equal(t, "secure_password", new_user.Password, "Password should be set to 'secure_password'")
	})

	t.Run("TestWithPasswordPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserBuilder().WithPassword("")
		}, "Setting an empty password should panic")
	})

	t.Run("TestWithEmail", func(t *testing.T) {
		new_user := NewUserBuilder().WithEmail("john@example.com")
		assert.Equal(t, "john@example.com", new_user.Email, "Email should be set to 'john@example.com'")
	})

	t.Run("TestWithEmailPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserBuilder().WithEmail("")
		}, "Setting an empty email should panic")

		assert.Panics(t, func() {
			NewUserBuilder().WithEmail("invalid_email")
		}, "Setting an invalid email should panic")
	})

	t.Run("TestWithAuthToken", func(t *testing.T) {
		token := "auth_token_123"
		new_user := NewUserBuilder().WithAuthToken(&token)
		assert.Equal(t, &token, new_user.AuthToken, "AuthToken should be set to 'auth_token_123'")
	})

	t.Run("TestWithEmailVerified", func(t *testing.T) {
		new_user := NewUserBuilder().WithEmailVerified(true)
		assert.True(t, new_user.EmailVerified, "EmailVerified should be set to true")

		new_user = NewUserBuilder().WithEmailVerified(false)
		assert.False(t, new_user.EmailVerified, "EmailVerified should be set to false")
	})

	t.Run("TestGetValidFields", func(t *testing.T) {
		token := "tok"
		user := NewUserBuilder().
			WithID(1).
			WithUsername("john").
			WithPassword("pass").
			WithEmail("john@example.com").
			WithAuthToken(&token).
			WithEmailVerified(true)
		user.CreatedAt = time.Now()

		fields := user.GetValidFields()
		assert.Contains(t, fields, "id")
		assert.Contains(t, fields, "username")
		assert.Contains(t, fields, "password")
		assert.Contains(t, fields, "email")
		assert.Contains(t, fields, "auth_token")
		assert.Contains(t, fields, "created_at")
		assert.Contains(t, fields, "email_verified")
	})

	t.Run("TestGet", func(t *testing.T) {
		token := "tok"
		user := NewUserBuilder().
			WithID(1).
			WithUsername("john").
			WithPassword("pass").
			WithEmail("john@example.com").
			WithAuthToken(&token).
			WithEmailVerified(true)
		user.CreatedAt = time.Now()

		val, ok := user.Get("id")
		assert.True(t, ok)
		assert.Equal(t, 1, val)

		val, ok = user.Get("username")
		assert.True(t, ok)
		assert.Equal(t, "john", val)

		val, ok = user.Get("password")
		assert.True(t, ok)
		assert.Equal(t, "pass", val)

		val, ok = user.Get("email")
		assert.True(t, ok)
		assert.Equal(t, "john@example.com", val)

		val, ok = user.Get("auth_token")
		assert.True(t, ok)
		assert.Equal(t, &token, val)

		_, ok = user.Get("created_at")
		assert.True(t, ok)

		val, ok = user.Get("email_verified")
		assert.True(t, ok)
		assert.Equal(t, true, val)

		val, ok = user.Get("not_a_field")
		assert.False(t, ok)
		assert.Nil(t, val)
	})

	t.Run("TestIsUserField", func(t *testing.T) {
		assert.True(t, IsUserField("id"))
		assert.True(t, IsUserField("username"))
		assert.True(t, IsUserField("password"))
		assert.True(t, IsUserField("email"))
		assert.True(t, IsUserField("auth_token"))
		assert.True(t, IsUserField("created_at"))
		assert.True(t, IsUserField("email_verified"))
		assert.False(t, IsUserField("not_a_field"))
	})

	t.Run("TestFieldsList", func(t *testing.T) {
		fields := FieldsList()
		expected := []string{"id", "username", "password", "email", "auth_token", "created_at", "email_verified"}
		assert.ElementsMatch(t, expected, fields)
	})

	t.Run("TestFieldsListWithout", func(t *testing.T) {
		fields := FieldsList("email", "auth_token")
		expected := []string{"id", "username", "password", "created_at", "email_verified"}
		assert.ElementsMatch(t, expected, fields)
	})
}
