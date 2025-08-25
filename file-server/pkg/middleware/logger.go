package middleware

import (
	"homeserver/file-server/internal/logger"
	"net/http"
)

func Logger(log logger.LoggerInt) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			log.Request(r)
		})
	}
}
