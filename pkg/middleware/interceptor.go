package middleware

import (
	"bytes"
	"fmt"
	"homeserver/pkg/parser"
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

		if reqCont.FullPageRender {
			w.Write(interceptor.body.Bytes())
			return
		}

		body_node, err := parser.ParseString(interceptor.body.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse response body: %v\n", err.Error())
			http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
			return
		}

		var modifiedContent *html.Node = nil

		if reqCont.ToRestorePage {
			modifiedContent, err = parser.MergePage(body_node)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to merge page content: %v\n", err.Error())
				http.Error(w, "Failed to merge page content", http.StatusInternalServerError)
				return
			}
		}

		if modifiedContent == nil {
			modifiedContent = body_node
		}

		if err := html.Render(w, modifiedContent); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write response: %v\n", err.Error())
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})
}
