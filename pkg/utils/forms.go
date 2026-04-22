package utils

import (
	"net/http"
)

type LoginForm struct {
	Identificator string `json:"identificator"`
	Password      string `json:"password"`
}

func LoginFormFromRequest(r *http.Request) LoginForm {
	return LoginForm{
		Identificator: r.FormValue("identificator"),
		Password:      r.FormValue("password"),
	}
}

type RegisterForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func RegisterFormFromRequest(r *http.Request) RegisterForm {
	return RegisterForm{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Confirm:  r.FormValue("confirm-password"),
	}
}
