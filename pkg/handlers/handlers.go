package handlers

import (
	"net/http"

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
