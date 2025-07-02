package handlers

import (
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"net/http"
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

func Login(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Handler] Login endpoint accessed")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("'user' form value", r.FormValue("user"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "../static/fragments/login.html")
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Handler] Register endpoint accessed")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`Error registering user`))

	// http.Redirect(w, r, "/", http.StatusSeeOther)
}
