package formatters

import "net/http"

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

	register_fields := []string{"username", "email", "password", "confirm"}
	for _, field := range register_fields {
		switch field {
		case "username":
			form.Username = r.FormValue(field)
		case "email":
			form.Email = r.FormValue(field)
		case "password":
			form.Password = r.FormValue(field)
		case "confirm":
			form.Confirm = r.FormValue(field)
		}
	}

	return form
}
