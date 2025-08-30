package middleware

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/internal/respfail"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Auth(conn *repositories.Repositories) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/static/") {
				next.ServeHTTP(w, r)
				return
			}

			sid := r.URL.Query().Get(constants.SIDQueryKey)

			expStr := r.URL.Query().Get(constants.ExpQueryKey)

			if expStr == "" {
				expStr = "0"
			}

			exp, err := strconv.ParseInt(expStr, 10, 64)
			if err != nil {
				respfail.ErrorResponse(w, r, "Unauthorized", http.StatusUnauthorized)
				return
			}

			pathParts := strings.Split(r.URL.Path, "/")[1:]
			if len(pathParts) != 2 {
				respfail.ErrorResponse(w, r, "Invalid request format", http.StatusBadRequest)
				return
			}

			entryType := pathParts[0]
			entryId, err := strconv.Atoi(pathParts[1])
			if err != nil {
				respfail.ErrorResponse(w, r, "Invalid request format", http.StatusBadRequest)
				return
			}

			trusted_source := ""
			owner := 0

			switch entryType {
			case "folders":
				fol := folschema.New().WithID(entryId)
				fol, err := conn.FolderRepo.GetFolder(r.Context(), fol, nil)
				if err != nil {
					respfail.ErrorResponse(w, r, "Invalid request format", http.StatusBadRequest)
					return
				}

				trusted_source = string(fol.TrustedSource)
				owner = fol.OwnerID

				r = r.WithContext(context.WithValue(r.Context(), constants.EntryKey, fol))
			case "files":
				file := fschema.New().WithID(entryId)
				f, err := conn.FileRepo.GetFile(r.Context(), file, nil)
				if err != nil {
					respfail.ErrorResponse(w, r, "Invalid request format", http.StatusBadRequest)
					return
				}

				trusted_source = string(f.TrustedSource)
				owner = f.OwnerID

				r = r.WithContext(context.WithValue(r.Context(), constants.EntryKey, f))
			default:
				respfail.ErrorResponse(w, r, "Invalid request format", http.StatusBadRequest)
				return
			}

			token := ""
			if authToken, _ := r.Cookie(constants.AuthCookieName); authToken != nil {
				token = authToken.Value
			}
			user, _ := conn.TokenRepo.QueryToken(r.Context(), token, true)

			fmt.Println("User ID:", user.ID)

			if token != "" && user.ID == owner {
				next.ServeHTTP(w, r)
				return
			}

			if sid == "" || trusted_source == "" || trusted_source != sid {
				respfail.ErrorResponse(w, r, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if exp < time.Now().Unix() {
				respfail.ErrorResponse(w, r, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
