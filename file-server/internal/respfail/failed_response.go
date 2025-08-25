package respfail

import (
	"context"
	"homeserver/constants"
	"homeserver/constants/fsconsts"
	"homeserver/file-server/internal/templates"
	"html/template"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, message string, status int) {
	r = r.WithContext(context.WithValue(r.Context(), fsconsts.StatusKey, status))

	tmpl, err := template.ParseFiles(constants.Rootd + "static/templates/load_error.html")
	if err != nil {
		http.Error(w, "Failed to load error template", http.StatusInternalServerError)
		return
	}

	data := templates.ErrorData{
		Title:   "Oops, we encountered an error",
		Message: message,
	}

	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render error template", http.StatusInternalServerError)
	}
}
