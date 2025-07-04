package globals

import (
	"fmt"
	"net/http"
	"strings"
)

type LoginForm struct {
	Identificator string `json:"identificator"`
	Password      string `json:"password"`
}

func LoginFormFromRequest(r *http.Request) (LoginForm, error) {
	var form LoginForm
	if err := r.ParseForm(); err != nil {
		return form, fmt.Errorf("failed to parse form: %w", err)
	}

	form.Identificator = r.FormValue("identificator")
	form.Password = r.FormValue("password")
	return form, nil
}

type RegisterForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func RegisterFormFromRequest(r *http.Request) (RegisterForm, error) {
	var form RegisterForm
	if err := r.ParseForm(); err != nil {
		return form, fmt.Errorf("failed to parse form: %w", err)
	}

	register_required_fields := []string{"username", "email", "password", "confirm-password"}
	missing_fields := []string{}
	for _, field := range register_required_fields {
		if r.FormValue(field) == "" {
			missing_fields = append(missing_fields, field)
		}

		switch field {
		case "username":
			form.Username = r.FormValue(field)
		case "email":
			form.Email = r.FormValue(field)
		case "password":
			form.Password = r.FormValue(field)
		case "confirm-password":
			form.Confirm = r.FormValue(field)
		}
	}

	if len(missing_fields) > 0 {
		return form, fmt.Errorf("missing required fields: %s", strings.Join(missing_fields, ", "))
	}

	if form.Password != form.Confirm {
		return form, fmt.Errorf("passwords do not match")
	}

	return form, nil
}
