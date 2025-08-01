package utils

import (
	"fmt"
	"homeserver/pkg/parser"
	"net/http"
)

type AlertGenerator func(title, iconName, message, status string) string

var (
	alertGenerator AlertGenerator = parser.NewAlert
)

const (
	AlertSuccess = "success"
	AlertError   = "error"
	AlertWarning = "warning"
	AlertInfo    = "info"
)

func ErrorResponse(w http.ResponseWriter, message string, status int) {
	fmt.Println("Error response:", message, "Status code:", status)

	alert := alertGenerator("Error", "x", message, AlertError)

	w.Header().Set("HX-Retarget", "#modal-content")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Push-Url", "false")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(alert))
}
