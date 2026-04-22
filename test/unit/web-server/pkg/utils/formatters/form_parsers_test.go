package formatters_test

import (
	"homeserver/constants/authconsts"
	"homeserver/pkg/utils"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginFormFromRequest(t *testing.T) {
	t.Run("Populates LoginForm correctly", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", nil)
		req.Form = map[string][]string{
			authconsts.IdentificatorField: {"testuser"},
			authconsts.PasswordField:      {"testpassword"},
		}

		form := formatters.LoginFormFromRequest(req)

		assert.Equal(t, "testuser", form.Identificator)
		assert.Equal(t, "testpassword", form.Password)
	})
}

func TestRegisterFormFromRequest(t *testing.T) {
	t.Run("Populates RegisterForm correctly", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/register", nil)
		req.Form = map[string][]string{
			authconsts.UsernameField:        {"testuser"},
			authconsts.EmailField:           {"testuser@example.com"},
			authconsts.PasswordField:        {"testpassword"},
			authconsts.ConfirmPasswordField: {"testpassword"},
		}

		form := formatters.RegisterFormFromRequest(req)

		assert.Equal(t, "testuser", form.Username)
		assert.Equal(t, "testuser@example.com", form.Email)
		assert.Equal(t, "testpassword", form.Password)
		assert.Equal(t, "testpassword", form.Confirm)
	})
}
