package handlers

import (
	"fmt"
	"homeserver/pkg/middleware"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Handler] Home page accessed")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "../static/index.html")
}

func Search(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Handler] Search endpoint accessed")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	results := `
		<p>Result 1</p>
		<p>Result 2</p>
		<p>Result 3</p>
	`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, "%s", results))
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Handler] Login page accessed")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "../static/fragments/login.html")
}
