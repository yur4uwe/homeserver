package utils

import (
	"homeserver/constants"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("HX-Trigger", AlertEventBuilder("Error", message, constants.Error))
	w.WriteHeader(status)
}
