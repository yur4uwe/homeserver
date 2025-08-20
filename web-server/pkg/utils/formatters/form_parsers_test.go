package formatters

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginFormFromRequest(t *testing.T) {
	t.Run("Populates LoginForm correctly", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", nil)
		req.Form = map[string][]string{
			"identificator": {"testuser"},
			"password":      {"testpassword"},
		}

		form := LoginFormFromRequest(req)

		assert.Equal(t, "testuser", form.Identificator)
		assert.Equal(t, "testpassword", form.Password)
	})
}

func TestRegisterFormFromRequest(t *testing.T) {
	t.Run("Populates RegisterForm correctly", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/register", nil)
		req.Form = map[string][]string{
			"username": {"testuser"},
			"email":    {"testuser@example.com"},
			"password": {"testpassword"},
			"confirm":  {"testpassword"},
		}

		form := RegisterFormFromRequest(req)

		assert.Equal(t, "testuser", form.Username)
		assert.Equal(t, "testuser@example.com", form.Email)
		assert.Equal(t, "testpassword", form.Password)
		assert.Equal(t, "testpassword", form.Confirm)
	})
}
