package middleware

import (
	"context"
	"homeserver/constants"
	"homeserver/db/repositories"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/file-server/internal/respfail"
	"net/http"
	"strconv"
	"strings"
)

func Auth(conn *repositories.Repositories) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sid := r.URL.Query().Get(constants.SIDQueryKey)

			pathParts := strings.Split(r.URL.Path, "/")[1:]
			if len(pathParts) != 2 {
				respfail.ErrorResponse(w, "Invalid request format", http.StatusBadRequest)
				return
			}

			entryType := pathParts[0]
			entryId, err := strconv.Atoi(pathParts[1])
			if err != nil {
				respfail.ErrorResponse(w, "Invalid request format", http.StatusBadRequest)
				return
			}

			trusted_source := ""
			owner := 0

			switch entryType {
			case "folders":
				fol := folschema.NewUserFolderBuilder().WithID(entryId)
				fol, err := conn.FolderRepo.GetFolder(r.Context(), fol)
				if err != nil {
					respfail.ErrorResponse(w, "Invalid request format", http.StatusBadRequest)
					return
				}

				trusted_source = string(fol.TrustedSource)
				owner = fol.OwnerID

				r = r.WithContext(context.WithValue(r.Context(), constants.EntryKey, fol))
			case "files":
				f, err := conn.FileRepo.GetFileByID(r.Context(), entryId)
				if err != nil {
					respfail.ErrorResponse(w, "Invalid request format", http.StatusBadRequest)
					return
				}

				trusted_source = string(f.TrustedSource)
				owner = f.OwnerID

				r = r.WithContext(context.WithValue(r.Context(), constants.EntryKey, f))
			default:
				respfail.ErrorResponse(w, "Invalid request format", http.StatusBadRequest)
				return
			}

			token := ""
			if authToken, _ := r.Cookie(constants.AuthCookieName); authToken != nil {
				token = authToken.Value
			}

			userId, err := conn.TokenRepo.CheckToken(r.Context(), token)
			if err != nil && !strings.Contains(err.Error(), "token") {
				respfail.ErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if (trusted_source == "" || trusted_source != sid) && owner != userId {
				respfail.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Wrapper(mux http.Handler, conn *repositories.Repositories) http.Handler {
	return Auth(conn)(mux)
}
