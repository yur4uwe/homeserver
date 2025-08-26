package formatters

import (
	"homeserver/constants/authconsts"
	"net/http"
)

type LoginForm struct {
	Identificator string `json:"identificator"`
	Password      string `json:"password"`
}

func LoginFormFromRequest(r *http.Request) LoginForm {
	var form LoginForm

	form.Identificator = r.FormValue("identificator")
	form.Password = r.FormValue("password")
	return form
}

type RegisterForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func RegisterFormFromRequest(r *http.Request) RegisterForm {
	var form RegisterForm

	register_fields := []string{
		authconsts.UsernameField,
		authconsts.EmailField,
		authconsts.PasswordField,
		authconsts.ConfirmPasswordField,
	}

	for _, field := range register_fields {
		switch field {
		case authconsts.UsernameField:
			form.Username = r.FormValue(field)
		case authconsts.EmailField:
			form.Email = r.FormValue(field)
		case authconsts.PasswordField:
			form.Password = r.FormValue(field)
		case authconsts.ConfirmPasswordField:
			form.Confirm = r.FormValue(field)
		}
	}

	return form
}
