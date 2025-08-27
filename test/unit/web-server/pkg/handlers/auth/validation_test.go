package auth_test

import (
	"homeserver/constants/authconsts"
	"homeserver/pkg/handlers/web/auth"
	"homeserver/pkg/utils/formatters"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginValidation(t *testing.T) {
	tests := []struct {
		name         string
		form         formatters.LoginForm
		wantValid    bool
		wantField    string
		wantContains string
	}{
		{
			name:         "empty identificator",
			form:         formatters.LoginForm{Identificator: "", Password: "pass"},
			wantValid:    false,
			wantField:    authconsts.IdentificatorField,
			wantContains: authconsts.IdentificatorError,
		},
		{
			name:         "empty password",
			form:         formatters.LoginForm{Identificator: "user", Password: ""},
			wantValid:    false,
			wantField:    authconsts.PasswordField,
			wantContains: authconsts.PasswordError,
		},
		{
			name:         "valid login",
			form:         formatters.LoginForm{Identificator: "user", Password: "pass"},
			wantValid:    true,
			wantField:    "",
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, field, html := auth.LoginValidation(tt.form)
			assert.Equal(t, tt.wantValid, valid)
			assert.Equal(t, tt.wantField, field)
			if tt.wantContains != "" {
				assert.True(t, strings.Contains(html, tt.wantContains))
			} else {
				assert.Equal(t, "", html)
			}
		})
	}
}

func TestRegisterValidation(t *testing.T) {
	tests := []struct {
		name         string
		form         formatters.RegisterForm
		wantValid    bool
		wantField    string
		wantContains string
	}{
		{
			name:         "empty username",
			form:         formatters.RegisterForm{Username: "", Email: "e", Password: "p", Confirm: "p"},
			wantValid:    false,
			wantField:    authconsts.UsernameField,
			wantContains: authconsts.UsernameError,
		},
		{
			name:         "empty email",
			form:         formatters.RegisterForm{Username: "u", Email: "", Password: "p", Confirm: "p"},
			wantValid:    false,
			wantField:    authconsts.EmailField,
			wantContains: authconsts.EmailError,
		},
		{
			name:         "empty password",
			form:         formatters.RegisterForm{Username: "u", Email: "e", Password: "", Confirm: "p"},
			wantValid:    false,
			wantField:    authconsts.PasswordField,
			wantContains: authconsts.PasswordError,
		},
		{
			name:         "empty confirm",
			form:         formatters.RegisterForm{Username: "u", Email: "e", Password: "p", Confirm: ""},
			wantValid:    false,
			wantField:    authconsts.ConfirmPasswordField,
			wantContains: authconsts.ConfirmPasswordError,
		},
		{
			name:         "password mismatch",
			form:         formatters.RegisterForm{Username: "u", Email: "e", Password: "p", Confirm: "q"},
			wantValid:    false,
			wantField:    authconsts.ConfirmPasswordField,
			wantContains: authconsts.PasswordMismatchError,
		},
		{
			name:         "valid register",
			form:         formatters.RegisterForm{Username: "u", Email: "e", Password: "p", Confirm: "p"},
			wantValid:    true,
			wantField:    "",
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, field, html := auth.RegisterValidation(tt.form)
			assert.Equal(t, tt.wantValid, valid)
			assert.Equal(t, tt.wantField, field)
			if tt.wantContains != "" {
				assert.True(t, strings.Contains(html, tt.wantContains))
			} else {
				assert.Equal(t, "", html)
			}
		})
	}
}
