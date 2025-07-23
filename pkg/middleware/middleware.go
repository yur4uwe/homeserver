package middleware

import (
	"context"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/parser"
	"homeserver/pkg/templates"
	"html/template"
	"net/http"
)

type contextKey string

const ReqContKey contextKey = "add_req_context"

type AdditionalRequestContext struct {
	LogHistory    *LogHistory
	ToRestorePage bool
	RerenderBody  bool
	AuthToken     string
	PageContent   *templates.LayoutPageData
}

func contextConstructor() *AdditionalRequestContext {
	return &AdditionalRequestContext{
		LogHistory:    &LogHistory{Entries: []string{}},
		ToRestorePage: false,
		RerenderBody:  false,
		AuthToken:     "",
		PageContent: &templates.LayoutPageData{
			Header: template.HTML(parser.MustAsString(parser.ParseFile("../static/templates/header.html"))),
			Footer: template.HTML(parser.MustAsString(parser.ParseFile("../static/fragments/footer.html"))),
		},
	}
}

func ToRestorePage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		if (r.Header.Get("HX-History-Restore-Request") == "true" || r.Header.Get("HX-Request") != "true") && globals.IsPage(r.URL.Path) {
			reqCont.ToRestorePage = true
			reqCont.LogHistory.Log("[Middleware] ToRestorePage set to true")

			w.Header().Set("Vary", "HX-History-Restore-Request, HX-Request")
		}

		next.ServeHTTP(w, r)
	})
}

func WithAdditionalRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont := contextConstructor()
		ctx := context.WithValue(r.Context(), ReqContKey, reqCont)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func MiddlewareWrapper(handler http.Handler, connection *db.Database) http.Handler {
	return WithAdditionalRequestContext(Logger(AuthMiddleware(connection)(ToRestorePage(Interceptor(handler)))))
}
