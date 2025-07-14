package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"homeserver/pkg/auth"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
)

func LoginHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in LoginHandler\n")
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Login endpoint accessed")

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse login form: %v\n", err.Error())
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		form := globals.LoginFormFromRequest(r)

		valid, target, fieldsetHTML := LoginValidation(form)
		if !valid {
			w.Header().Set("HX-Retarget", "#"+target+"-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			w.Write([]byte(fieldsetHTML))
			return
		}

		userId, email_verified, err := connection.QueryUserOnLogin(r.Context(), form)
		if err != nil {
			if err.Error() == "no rows in result set" {
				fmt.Fprintf(os.Stderr, "Login failed for user %s: %v\n", form.Identificator, err.Error())

				w.Header().Set("HX-Retarget", "#error-resp")
				w.Header().Set("HX-Reswap", "innerHTML")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized: Invalid credentials"))
				return
			}
			fmt.Fprintf(os.Stderr, "Failed to query user on login: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if !email_verified {
			fmt.Fprintf(os.Stderr, "User %d has not verified their email\n", userId)
			w.Header().Set("HX-Retarget", "#error-resp")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden: Email not verified"))
			return
		}

		newToken, err := auth.NewToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate new token: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = connection.CreateToken(r.Context(), userId, newToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create token: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authTokenHS",
			Value:    newToken,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		index, err := parser.ParseFile("../static/index.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse index HTML file: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		body := parser.GetNodeByTagName(index, "body")
		if body == nil {
			fmt.Fprintf(os.Stderr, "Failed to find body tag in index HTML\n")
			http.Error(w, "Internal Server Error: Body tag not found", http.StatusInternalServerError)
			return
		}

		welcomeNode, err := parser.ParseFile("../static/fragments/welcome_logged_in.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse welcome HTML file: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		parser.InsertNodeIntoPage(body, welcomeNode, "main-content")

		w.Header().Set("HX-Push-Url", "/")
		w.WriteHeader(http.StatusOK)

	}
}

func RegisterHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in RegisterHandler\n")
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse registration form: %v\n", err.Error())
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		reqCont.LogHistory.Log("[Handler] Register endpoint accessed")

		form := globals.RegisterFormFromRequest(r)

		valid, target, fieldsetHTML := RegisterValidation(form)
		if !valid {
			w.Header().Set("HX-Retarget", "#"+target+"-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			w.Write([]byte(fieldsetHTML))
			return
		}

		userId, err := connection.CreateUser(r.Context(), form.Username, form.Email, form.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		verification_token, err := auth.NewToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate verification token: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = connection.CreateToken(r.Context(), userId, verification_token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create verification token: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		go sendVerificationEmail(form, verification_token)

		alert := parser.NewAlert("You have registered successfully. Now verify your email address.", parser.AlertSuccess)

		index, err := parser.ParseFile("../static/index.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse index HTML file: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		body := parser.GetNodeByTagName(index, "body")
		if body == nil {
			fmt.Fprintf(os.Stderr, "Failed to find body tag in index HTML\n")
			http.Error(w, "Internal Server Error: Body tag not found", http.StatusInternalServerError)
			return
		}

		welcomeNode, err := parser.ParseFile("../static/fragments/welcome.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse welcome HTML file: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		parser.InsertNodeIntoPage(body, welcomeNode, "main-content")
		parser.InsertNodeIntoPage(body, alert, "modal-content")

		reqCont.FullPageRender = body

		w.Header().Set("HX-Push-Url", "/")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
}
