package handlers

import (
	"homeserver/pkg/middleware"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	err := middleware.AddLog(r, "[Hello Handler] Hello handler called")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
