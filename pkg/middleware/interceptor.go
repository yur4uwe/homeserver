package middleware

import (
	"bytes"
	"fmt"
	"homeserver/pkg/parser"
	"html/template"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

type responseInterceptor struct {
	statusCode int
	body       *bytes.Buffer
	http.ResponseWriter
}

func InterceptorConstructor(w http.ResponseWriter) *responseInterceptor {
	return &responseInterceptor{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           new(bytes.Buffer),
	}
}

func (ri *responseInterceptor) WriteHeader(code int) {
	ri.statusCode = code
}

func (ri *responseInterceptor) Write(b []byte) (int, error) {
	return ri.body.Write(b)
}

func Interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in interceptor middleware\n")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		interceptor := InterceptorConstructor(w)

		next.ServeHTTP(interceptor, r)

		if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
			if _, err := w.Write(interceptor.body.Bytes()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write response: %v\n", err.Error())
			}
			return
		}

		if reqCont.ToRestorePage {
			reqCont.LogHistory.Log("[Interceptor] ToRestorePage is true, restoring page")
			tmpl := template.Must(template.ParseFiles("../static/index.html", "../static/layout.html"))

			reqCont.PageContent.Content = template.HTML(interceptor.body.String())

			if err := tmpl.ExecuteTemplate(w, "index.html", reqCont.PageContent); err != nil {
				reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "Interceptor:89")
				http.Error(w, "Failed to render page", http.StatusInternalServerError)
			}
			return
		}

		if reqCont.RerenderBody {
			reqCont.LogHistory.Log("[Interceptor] RerenderBody is true, rendering body")
			w.Header().Set("HX-Retarget", "body")
			w.Header().Set("HX-Reswap", "innerHTML")

			tmpl := template.Must(template.ParseFiles("../static/layout.html"))

			if err := tmpl.Execute(w, reqCont.PageContent); err != nil {
				reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "Interceptor:60")
				http.Error(w, "Failed to render page", http.StatusInternalServerError)
			}
			return
		}

		body_node, err := parser.ParseString(interceptor.body.String())
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse response body: "+err.Error()+"\nBody: "+interceptor.body.String(), "Interceptor:65")
			http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
			return
		}

		if err := html.Render(w, body_node); err != nil {
			reqCont.LogHistory.ErrLog("Failed to render HTML response: "+err.Error(), "Interceptor:86")
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})
}
