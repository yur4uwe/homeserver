package errresp

import (
	"fmt"
	"net/http"
)

const (
	AlertSuccess = "success"
	AlertError   = "error"
	AlertWarning = "warning"
	AlertInfo    = "info"
)

func ErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("HX-Trigger", fmt.Sprintf(`{"raise-alert": {"title": "Error", "value": "%s", "status": "%s", "icon": "x"}}`, message, AlertError))
	w.WriteHeader(status)
}
