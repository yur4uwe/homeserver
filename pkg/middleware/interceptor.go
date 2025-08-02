package middleware

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func Interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in interceptor middleware\n")
			ErrorResponse(w, "Internal Server Error", http.StatusInternalServerError) // Use the global variable for mocking
			return
		}

		interceptor := InterceptorConstructor(w)

		next.ServeHTTP(interceptor, r)

		reqCont.LogHistory.Log("[Interceptor] Response status code: " + fmt.Sprint(interceptor.statusCode))

		if interceptor.statusCode >= 400 {
			w.WriteHeader(interceptor.statusCode)
		}

		if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
			if _, err := w.Write(interceptor.body.Bytes()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write response: %v\n", err.Error())
			}
			return
		}

		wr := &WriterWithIcons{ResponseWriter: w}

		if reqCont.ToRestorePage {
			reqCont.LogHistory.Log("[Interceptor] ToRestorePage is true, restoring page")
			tmpl := template.Must(ParseFiles("../static/index.html", "../static/layout.html")) // Use the global variable for mocking

			if interceptor.statusCode < 400 {
				reqCont.PageContent.Content = template.HTML(interceptor.body.String())
			} else {
				reqCont.PageContent.ModalContent = template.HTML(interceptor.body.String())
			}

			if err := tmpl.ExecuteTemplate(wr, "index.html", reqCont.PageContent); err != nil {
				reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "Interceptor:89")
				ErrorResponse(w, "Failed to render page", http.StatusInternalServerError) // Use the global variable for mocking
			}
			return
		}

		if reqCont.RerenderBody {
			reqCont.LogHistory.Log("[Interceptor] RerenderBody is true, rendering body")
			w.Header().Set("HX-Retarget", "body")
			w.Header().Set("HX-Reswap", "innerHTML")

			tmpl := template.Must(ParseFiles("../static/layout.html")) // Use the global variable for mocking

			if interceptor.statusCode < 400 {
				reqCont.PageContent.Content = template.HTML(interceptor.body.String())
			} else {
				reqCont.PageContent.ModalContent = template.HTML(interceptor.body.String())
			}

			if err := tmpl.Execute(wr, reqCont.PageContent); err != nil {
				reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "Interceptor:60")
				ErrorResponse(w, "Failed to render page", http.StatusInternalServerError) // Use the global variable for mocking
			}
			return
		}

		wr.Write(interceptor.body.Bytes())
	})
}
