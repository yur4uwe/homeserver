package handlers

import (
	"net/http"
	"time"

	"homeserver/pkg/auth"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
)

func PageServer(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Handler] Page server accessed for page: "+r.URL.Path)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !globals.IsPage(r.URL.Path) {
		http.NotFound(w, r)
		return
	}

	file_to_serve := "../static/fragments/"
	if r.URL.Path == "/" {
		file_to_serve += "welcome"
	} else {
		file_to_serve += r.URL.Path[1:]
	}
	file_to_serve += ".html"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, file_to_serve)
}

func LoginHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := middleware.AddLog(r, "[Handler] Login endpoint accessed")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		form, err := globals.LoginFormFromRequest(r)
		if err != nil {
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := connection.QueryUserOnLogin(r.Context(), form)
		if err != nil {
			if err.Error() == "no rows in result set" {
				http.Error(w, "Unauthorized: Invalid credentials", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newToken := auth.NewToken(userId)

		http.SetCookie(w, &http.Cookie{
			Name:     "authTokenHS",
			Value:    newToken,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RegisterHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := middleware.AddLog(r, "[Handler] Register endpoint accessed")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		form, err := globals.RegisterFormFromRequest(r)
		if err != nil {
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		if form.Password != form.Confirm {
			http.Error(w, "Bad Request: Passwords do not match", http.StatusBadRequest)
			return
		}

		_, err = connection.CreateUser(r.Context(), form.Username, form.Email, form.Password)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		alert := `
			<div class="alert alert-success" role="alert">
				<strong>Success!</strong> You have registered successfully. Now log in
			</div>
		`

		response, err := parser.InsertNodeIntoPage("../static/fragments/welcome.html", alert, "")
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
