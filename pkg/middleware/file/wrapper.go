package middleware

import (
	"homeserver/db/repositories"
	"homeserver/internal/logger"
	"net/http"
)

func Wrapper(mux http.Handler, conn *repositories.Repositories, log logger.LoggerInt) http.Handler {
	return Logger(log)(Auth(conn)(mux))
}
