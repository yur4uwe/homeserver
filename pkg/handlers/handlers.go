package handlers

import (
	"fmt"
	"net/http"

	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
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

		fmt.Println("Login form received:", form)

		id, err := connection.QueryUserOnLogin(r.Context(), form)
		if err != nil {
			if err.Error() == "no rows in result set" {
				http.Error(w, "Unauthorized: Invalid credentials", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("User ID:", id)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "../static/fragments/login.html")
	}
}

func RegisterHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := middleware.AddLog(r, "[Handler] Register endpoint accessed")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Println("'user' form value", r.FormValue("user"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "../static/fragments/register.html")
	}
}
