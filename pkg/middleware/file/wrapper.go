package middleware

import (
	"homeserver/db"
	"homeserver/internal/logger"
	"net/http"
)

func Wrapper(mux http.Handler, conn *db.Database, log logger.LoggerInt) http.Handler {
	return Logger(log)(Auth(conn)(mux))
}
