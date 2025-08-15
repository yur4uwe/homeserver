package auth

import (
	"homeserver/pkg/globals"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginValidation(t *testing.T) {
	tests := []struct {
		name         string
		form         globals.LoginForm
		wantValid    bool
		wantField    string
		wantContains string
	}{
		{
			name:         "empty identificator",
			form:         globals.LoginForm{Identificator: "", Password: "pass"},
			wantValid:    false,
			wantField:    IdentificatorField,
			wantContains: IdentificatorError,
		},
		{
			name:         "empty password",
			form:         globals.LoginForm{Identificator: "user", Password: ""},
			wantValid:    false,
			wantField:    PasswordField,
			wantContains: PasswordError,
		},
		{
			name:         "valid login",
			form:         globals.LoginForm{Identificator: "user", Password: "pass"},
			wantValid:    true,
			wantField:    "",
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, field, html := LoginValidation(tt.form)
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
		form         globals.RegisterForm
		wantValid    bool
		wantField    string
		wantContains string
	}{
		{
			name:         "empty username",
			form:         globals.RegisterForm{Username: "", Email: "e", Password: "p", Confirm: "p"},
			wantValid:    false,
			wantField:    UsernameField,
			wantContains: UsernameError,
		},
		{
			name:         "empty email",
			form:         globals.RegisterForm{Username: "u", Email: "", Password: "p", Confirm: "p"},
			wantValid:    false,
			wantField:    EmailField,
			wantContains: EmailError,
		},
		{
			name:         "empty password",
			form:         globals.RegisterForm{Username: "u", Email: "e", Password: "", Confirm: "p"},
			wantValid:    false,
			wantField:    PasswordField,
			wantContains: PasswordError,
		},
		{
			name:         "empty confirm",
			form:         globals.RegisterForm{Username: "u", Email: "e", Password: "p", Confirm: ""},
			wantValid:    false,
			wantField:    ConfirmPasswordField,
			wantContains: ConfirmPasswordError,
		},
		{
			name:         "password mismatch",
			form:         globals.RegisterForm{Username: "u", Email: "e", Password: "p", Confirm: "q"},
			wantValid:    false,
			wantField:    ConfirmPasswordField,
			wantContains: PasswordMismatchError,
		},
		{
			name:         "valid register",
			form:         globals.RegisterForm{Username: "u", Email: "e", Password: "p", Confirm: "p"},
			wantValid:    true,
			wantField:    "",
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, field, html := RegisterValidation(tt.form)
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
