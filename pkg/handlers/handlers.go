package handlers

import (
	"fmt"
	"net/http"
	"os"

	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
)

func PageServer(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		fmt.Fprintf(os.Stderr, "Missing request context in PageServer handler\n")
		http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Page server accessed for page: " + r.URL.Path)

	if !globals.IsPage(r.URL.Path) {
		reqCont.LogHistory.Log("Path not considered a page: " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/" {
		WelcomePageHandler(w, r)
		return
	}

	file_to_serve := "../static/fragments/" + r.URL.Path[1:] + ".html"

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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		http.Error(w, "Invalid request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("Logout request received")

	// Clear the session or authentication token
	reqCont.AuthToken = ""

	// Redirect to the login page or home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
