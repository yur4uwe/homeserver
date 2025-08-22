package routes

import (
	"homeserver/constants"
	"homeserver/db/repositories"
	"homeserver/web-server/pkg/handlers"
	"homeserver/web-server/pkg/handlers/auth"
	"homeserver/web-server/pkg/middleware"
	"net/http"
)

func NewRouter(repos *repositories.Repositories) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(constants.Rootd + "static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/verification", auth.VerificationHandler(repos))
	mux.HandleFunc("/upload", handlers.UploadPageHandler)
	mux.HandleFunc("/settings", handlers.SettingsPageHandler)
	mux.HandleFunc("/search/", handlers.SearchRouter(repos))
	mux.HandleFunc("/api/v1/upload", handlers.UploadRouter(repos))
	mux.HandleFunc("/api/v1/login", auth.LoginHandler(repos))
	mux.HandleFunc("/api/v1/register", auth.RegisterHandler(repos))
	mux.HandleFunc("/api/v1/logout", auth.LogoutHandler(repos))
	mux.HandleFunc("/api/v1/file-actions/", handlers.FileActionsRouter(repos))
	mux.HandleFunc("/api/v1/account/settings", handlers.AccountSettingsRouter(repos))
	mux.HandleFunc("/", handlers.StaticPageHandler)

	handler := middleware.MiddlewareWrapper(mux, repos)

	return handler
}
