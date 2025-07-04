package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"homeserver/pkg/auth"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"

	"golang.org/x/net/html"
)

func LoginHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := middleware.AddLog(r, "[Handler] Login endpoint accessed")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		form, err := globals.LoginFormFromRequest(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse login form: %v\n", err.Error())
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := connection.QueryUserOnLogin(r.Context(), form)
		if err != nil {
			if err.Error() == "no rows in result set" {
				fmt.Fprintf(os.Stderr, "Login failed for user %s: %v\n", form.Identificator, err.Error())
				http.Error(w, "Unauthorized: Invalid credentials", http.StatusUnauthorized)
				return
			}
			fmt.Fprintf(os.Stderr, "Failed to query user on login: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newToken := auth.NewToken(userId)

		err = connection.CreateToken(r.Context(), userId, newToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create token: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

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
			fmt.Fprintf(os.Stderr, "Failed to parse registration form: %v\n", err.Error())
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		if form.Password != form.Confirm {
			fmt.Fprintf(os.Stderr, "Passwords do not match: %s != %s\n", form.Password, form.Confirm)
			http.Error(w, "Bad Request: Passwords do not match", http.StatusBadRequest)
			return
		}

		_, err = connection.CreateUser(r.Context(), form.Username, form.Email, form.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		alert, err := parser.ParseString(`
			<div class="alert alert-success" role="alert">
				<strong>Success!</strong> You have registered successfully. Now log in
			</div>
		`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse alert HTML: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// w.Header().Set("HX-Trigger", "alert:show")
		w.Header().Set("Location", "/")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusSeeOther)
		html.Render(w, alert)
	}
}
