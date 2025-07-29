package utils

import (
	"homeserver/pkg/parser"
	"net/http"
)

const (
	AlertSuccess = "success"
	AlertError   = "error"
	AlertWarning = "warning"
	AlertInfo    = "info"
)

func ErrorResponse(w http.ResponseWriter, message string, status int) {
	alert := parser.NewAlert("Error", "error", message, AlertError)

	w.Header().Set("HX-Retarget", "#modal-content")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(alert))
}
