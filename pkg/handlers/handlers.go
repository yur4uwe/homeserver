package handlers

import (
	"fmt"
	"net/http"
	"os"

	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
)

func PageServer(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Handler] Page server accessed for page: "+r.URL.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to log request: %v\n", err.Error())
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

	content, err := os.ReadFile(file_to_serve)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file %s: %v\n", file_to_serve, err.Error())
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write reponse: %v\n", err.Error())
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func NothingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}
