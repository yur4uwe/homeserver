package utils

import (
	"fmt"
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
	fmt.Println("Error response:", message, "Status code:", status)

	alert := parser.NewAlert("Error", "x", message, AlertError)

	w.Header().Set("HX-Retarget", "#modal-content")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Push-Url", "false")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(alert))
}
