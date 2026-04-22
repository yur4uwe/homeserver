package routes

import (
	"homeserver/constants"
	"homeserver/db"
	"homeserver/internal/logger"
	filehandlers "homeserver/pkg/handlers/file"
	webhandlers "homeserver/pkg/handlers/web"
	"homeserver/pkg/handlers/web/auth"
	filemiddleware "homeserver/pkg/middleware/file"
	webmiddleware "homeserver/pkg/middleware/web"
	"net/http"
)

func NewWebRouter(repos *db.Database) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(constants.Rootd + "static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/verification", auth.VerificationHandler(repos))
	mux.HandleFunc("/upload", webhandlers.UploadPageHandler)
	mux.HandleFunc("/settings", webhandlers.SettingsPageHandler)
	mux.HandleFunc("/search/", webhandlers.SearchRouter(repos))
	mux.HandleFunc("/api/v1/upload", webhandlers.UploadRouter(repos))
	mux.HandleFunc("/api/v1/login", auth.LoginHandler(repos))
	mux.HandleFunc("/api/v1/register", auth.RegisterHandler(repos))
	mux.HandleFunc("/api/v1/logout", auth.LogoutHandler(repos))
	mux.HandleFunc("/api/v1/file-actions/", webhandlers.FileActionsRouter(repos))
	mux.HandleFunc("/api/v1/account/settings", webhandlers.AccountSettingsRouter(repos))
	mux.HandleFunc("/", webhandlers.StaticPageHandler)

	handler := webmiddleware.MiddlewareWrapper(mux, repos)

	return handler
}

func NewFileRouter(repos *db.Database) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(constants.Rootd + "static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/files/", filehandlers.FileHandler)
	mux.HandleFunc("/folders/", filehandlers.FolderHandler)

	handler := filemiddleware.Wrapper(mux, repos, &logger.Logger{})

	return handler
}
