package middleware

import (
	"context"
	"homeserver/db/repositories"
	"homeserver/web-server/pkg/templates"
	"homeserver/web-server/pkg/utils"
	"net/http"
)

var (
	IsPage = utils.IsPage
)

func contextConstructor() *AdditionalRequestContext {
	return &AdditionalRequestContext{
		LogHistory:    &LogHistory{Entries: []string{}},
		ToRestorePage: false,
		RerenderBody:  false,
		AuthToken:     "",
		Requester:     nil,
		PageContent:   &templates.LayoutPageData{},
	}
}

func ToRestorePage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		if (r.Header.Get("HX-History-Restore-Request") == "true" || r.Header.Get("HX-Request") != "true") && IsPage(r.URL.Path) {
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

func MiddlewareWrapper(handler http.Handler, repos *repositories.Repositories) http.Handler {
	return WithAdditionalRequestContext(Logger(AuthMiddleware(repos)(ToRestorePage(Interceptor(handler)))))
}
