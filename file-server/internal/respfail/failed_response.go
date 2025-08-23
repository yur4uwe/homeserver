package respfail

import "net/http"

func ErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Disposition", `attachment; filename="error.txt"`)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Download failed: " + message))
}
