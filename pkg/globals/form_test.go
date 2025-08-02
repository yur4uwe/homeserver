package globals

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginFormFromRequest(t *testing.T) {
	t.Run("Populates LoginForm correctly", func(t *testing.T) {
		// Create a mock HTTP request with form values
		req := httptest.NewRequest("POST", "/login", nil)
		req.Form = map[string][]string{
			"identificator": {"testuser"},
			"password":      {"testpassword"},
		}

		// Call LoginFormFromRequest
		form := LoginFormFromRequest(req)

		// Assert that the form is populated correctly
		assert.Equal(t, "testuser", form.Identificator)
		assert.Equal(t, "testpassword", form.Password)
	})
}

func TestRegisterFormFromRequest(t *testing.T) {
	t.Run("Populates RegisterForm correctly", func(t *testing.T) {
		// Create a mock HTTP request with form values
		req := httptest.NewRequest("POST", "/register", nil)
		req.Form = map[string][]string{
			"username":         {"testuser"},
			"email":            {"testuser@example.com"},
			"password":         {"testpassword"},
			"confirm-password": {"testpassword"},
		}

		// Call RegisterFormFromRequest
		form := RegisterFormFromRequest(req)

		// Assert that the form is populated correctly
		assert.Equal(t, "testuser", form.Username)
		assert.Equal(t, "testuser@example.com", form.Email)
		assert.Equal(t, "testpassword", form.Password)
		assert.Equal(t, "testpassword", form.Confirm)
	})
}
