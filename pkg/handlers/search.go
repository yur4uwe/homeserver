package handlers

import (
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"net/http"
)

func SearchHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Search API accessed")

		query := r.URL.Query().Get("query")
		directory := r.URL.Query().Get("directory")

		entries, err := connection.SearchFiles(r.Context(), query, reqCont.AuthToken)
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			http.Error(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
			return
		}

		data := &templates.SearchPageData{
			Query:     query,
			Directory: directory,
			Entries:   entries,
		}

		SearchPageHandler(data)(w, r)
	}
}
