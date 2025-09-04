package errresp

import (
	"homeserver/constants"
	"homeserver/pkg/utils/formatters"
	"net/http"
)

const (
	AlertSuccess = "success"
	AlertError   = "error"
	AlertWarning = "warning"
	AlertInfo    = "info"
)

func ErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("HX-Trigger", formatters.AlertEventBuilder("Error", message, constants.Error))
	w.WriteHeader(status)
}
